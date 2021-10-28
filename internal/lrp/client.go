package lrp

import (
	"encoding/binary"
	"errors"
	"lrp/internal/common"
	nt "lrp/internal/conn"
	"net"
	"time"
)

type Client struct {
	id       []byte
	token    uint32
	verbose  bool
	presult  chan []byte
	Conn     nt.Conn
	TrBucket *common.Bucket
}

func NewClient(tk uint32, vb bool) *Client {
	if tk == 0 {
		log.Error("input token error")
		return nil
	}
	return &Client{
		token:    uint32(tk),
		verbose:  vb,
		presult:  make(chan []byte),
		TrBucket: common.NewBucket(1024),
	}
}

func (c *Client) Run(server, proto string) error {
	nt, err := nt.NewConn(proto)
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
					log.Error("recive data error: ", err)
					log.Error("exit ...")
					return
				} else {
					switch pl[0] {
					case 1:
						switch pl[1] {
						case 3:
							c.presult <- pl[2:]
						default:
							log.Warn("Unsupported replay type")
						}
					case 2:
						reply, isCreateOk := append(append([]byte{1, 2}, pl[1:13]...), pl[25:37]...), byte(1)
						dest := common.AddrByteToString(pl[37:])
						if err := c.NewTransport(dest, pl[1:13], pl[13:25]); err != nil {
							log.Warn("create transport failed", err)
							isCreateOk = byte(0)
						}
						if _, err := EncodeSend(c.Conn, append(reply, isCreateOk)); err != nil {
							log.Warn("send to Server failed(on createTr)", err)
						}
					case 3:
						if tr := c.TrBucket.Get(common.XidToString(pl[1:13])); tr != nil {
							if err := tr.(*Transport).Write(pl[13:]); err != nil {
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

func (c *Client) AddTempProxy(dest string) (uint16, error) {
	if addr, err := common.AddrStringToByte(dest, "tcp"); err != nil {
		return 0, err
	} else {
		req := append([]byte{4}, addr...)
		if _, err := EncodeSend(c.Conn, req); err != nil {
			return 0, err
		}
		for {
			select {
			case <-time.After(time.Second * 5):
				return 0, errors.New("wait server reply timeout")
			case res := <-c.presult:
				if res[0] != 1 {
					return 0, errors.New("request temp proxy failed")
				}
				return binary.BigEndian.Uint16(res[1:]), nil
			}
		}
	}
}

func (c *Client) NewTransport(dest string, pid, tid []byte) error {
	if conn, err := net.Dial("tcp", dest); err != nil {
		return err
	} else {
		tr := NewTransport(false, tid, pid, c.Conn, conn)
		c.TrBucket.Set(common.XidToString(tid), tr)
		go tr.Serve()
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
