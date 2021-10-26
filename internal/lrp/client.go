package lrp

import (
	"encoding/binary"
	"errors"
	"lrp/internal/common"
	nt "lrp/internal/conn"
	"net"
	"strconv"
	"time"
)

type Client struct {
	id       []byte
	token    uint32
	verbose  bool
	Conn     nt.Conn
	TrBucket *common.Bucket
}

func NewClient(tk string, vb bool) *Client {
	if tk == "" || len(tk) != 6 {
		log.Error("input token error")
		return nil
	}
	ti, err := strconv.Atoi(tk)
	if err != nil {
		log.Error("decode token error")
		return nil
	}
	return &Client{
		token:    uint32(ti),
		verbose:  vb,
		TrBucket: common.NewBucket(1024),
	}
}

func (c *Client) Run(server, dest string) error {
	tmpResult := make(chan []byte)
	nt, err := nt.NewConn("tcp")
	if err != nil {
		return err
	}
	if c.Conn, err = nt.Dial(server); err != nil {
		return err
	} else {
		if c.id, err = c.auth(); err != nil {
			return err
		}
		go func() {
			defer c.Conn.Close()
			for {
				if pl, err := DecodeReceive(c.Conn); err != nil {
					log.Error("recive data error", err)
				} else {
					switch pl[0] {
					case 1:
						//临时代理请求结果
						if pl[1] == 3 {
							tmpResult <- pl[2:]
						}
						log.Warn("Unsupported result type", pl)
					case 2:
						//新连接请求
						reply := append([]byte{1, 2}, pl[1:9]...)
						isCreateOk := byte(1)
						dest := common.AddrByteToString(pl[17:])
						if err := c.NewTransport(dest, pl[1:9], pl[9:17]); err != nil {
							log.Warn("create transport failed", err)
							isCreateOk = byte(0)
						}
						if _, err := EncodeSend(c.Conn, append(reply, isCreateOk)); err != nil {
							log.Warn("send to Server failed(on createTr)", err)
						}
					case 3:
						if tr := c.TrBucket.Get(common.XidToString(pl[1:9])); tr != nil {
							if err := tr.(*Transport).Write(pl[9:]); err != nil {
								log.Warn("write data to dest err", err)
							}
						}
					default:
						log.Error("Unsupported protocol")
						return
					}
				}
			}
		}()
		return nil
	}
}

func (c *Client) NewTransport(dest string, pid, tid []byte) error {
	if conn, err := net.Dial("tcp", dest); err != nil {
		return err
	} else {
		tr := NewTransport(false, tid, pid, c.Conn, conn)
		c.TrBucket.Set(common.XidToString(tid), tr)
		tr.Serve()
		return nil
	}
}

func (c *Client) auth() ([]byte, error) {
	tk := make([]byte, 4)
	binary.BigEndian.PutUint32(tk, c.token)
	req := append([]byte{0}, tk...)
	if _, err := EncodeSend(c.Conn, req); err != nil {
		return nil, err
	} else {
		for {
			select {
			case <-time.After(time.Second * 5):
				return nil, errors.New("wait auth reply timeout")
			default:
				if data, err := DecodeReceive(c.Conn); err != nil {
					return nil, err
				} else {
					if data[0] == 1 && data[1] == 1 {
						if data[2] == 1 {
							log.Info("server auth successful , clientID: " + common.XidToString(data[3:]))
							return data[3:], nil
						}
						return nil, errors.New("server auth failed")
					}
					log.Error("not support data", data)
					return nil, errors.New("not support cmd")
				}
			}
		}
	}
}
