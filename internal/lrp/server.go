package lrp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"lrp/internal/common"
	nt "lrp/internal/conn"
	"strconv"

	"github.com/rs/xid"
)

type Server struct {
	token      uint32
	proto      string
	clients    *common.Bucket
	version    string
	ListenPort string
	ExternalIp string
}

func NewServer() *Server {
	return &Server{
		version: version,
		clients: common.NewBucket(1024),
	}
}

func (s *Server) Run(sp, proto string, tk uint32) (uint32, bool) {
	var err error
	defer func() {
		if err != nil {
			log.Warn(err.Error())
		}
	}()

	if tk == 0 {
		if newTk := common.GenerateServerToken(); newTk == "" {
			err = errors.New("generate token failed")
			return 0, false
		} else {
			tk, _ := strconv.Atoi(newTk)
			s.token = uint32(tk)
		}
	} else {
		s.token = tk
	}

	if s.ExternalIp, err = common.GetExternalIp(); err != nil {
		return 0, false
	}

	ns, err := nt.NewConn(proto)
	if err != nil {
		return 0, false
	}

	ln, err := ns.Listen(":" + sp)
	if err != nil {
		return 0, false
	}

	go func() {
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Warn("accept new client error: ", err)
			} else {
				go s.handleClient(conn)
			}
		}
	}()
	s.proto, s.ListenPort = proto, ln.Info()
	return s.token, true
}

func (s *Server) handleClient(conn nt.Conn) {
	if tk, err := DecodeReceive(conn); err != nil || binary.BigEndian.Uint32(tk[1:]) != s.token {
		log.Warn("auth client faild ", binary.BigEndian.Uint32(tk[1:]))
		EncodeSend(conn, []byte{1, 1, 0})
		return
	}
	id := xid.New()
	sc := newSClient(conn, id.String())
	if _, err := EncodeSend(conn, append([]byte{1, 1, 1}, id.Bytes()...)); err != nil {
		log.Warn("send auth reply err ", err)
		return
	}
	s.clients.Set(id.String(), sc)
	defer func() {
		conn.Close()
		s.clients.Remove(id.String())
	}()
	sc.Serve()
}

func (s *Server) AddProxy(cid, dest, mark, listenPort string) error {
	if client := s.clients.Get(cid); client != nil {
		if _, err := client.(*SClient).AddProxy(dest, mark, listenPort, false); err != nil {
			return err
		}
		return nil
	}
	return errors.New("cant find client id")
}

func (s *Server) DelProxy(cid, pid string) error {
	if client := s.clients.Get(cid); client != nil {
		if err := client.(*SClient).DelProxy(pid); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) CheckToken(value uint32) bool {
	return s.token == value
}

func (s *Server) GetServerInfo() ServerInfo {
	si := ServerInfo{
		Version:     s.version,
		Protocol:    s.proto,
		ProxyNum:    0,
		TProxyNum:   0,
		ClientNum:   0,
		ExternalIp:  s.ExternalIp,
		ServerPort:  s.ListenPort,
		ClientInfos: make([]ClientInfo, 0),
	}
	for _, c := range s.clients.GetAll() {
		si.ClientNum++
		client := c.(*SClient)
		ci := ClientInfo{
			Id:         client.id,
			Mark:       "",
			Online:     client.online,
			ProxyInfos: make([]ProxyInfo, 0),
		}
		for _, p := range client.proxyBucket.GetAll() {
			ps := p.(*ProxyServer)
			if !ps.Temp {
				si.ProxyNum++
			} else {
				si.TProxyNum++
			}
			pi := ProxyInfo{
				Id:      common.XidToString(ps.Id),
				Info:    fmt.Sprintf("%s:%d => %s", si.ExternalIp, ps.ListenPort, common.AddrByteToString(ps.DestAddr)),
				Mark:    "",
				Status:  ps.Status,
				ConnNum: len(ps.TransportBucket.GetAll()),
			}
			ci.ProxyInfos = append(ci.ProxyInfos, pi)
		}
		si.ClientInfos = append(si.ClientInfos, ci)
	}
	return si
}

type SClient struct {
	id          string
	conn        nt.Conn
	online      bool
	proxyBucket *common.Bucket
}

func newSClient(conn nt.Conn, id string) *SClient {
	return &SClient{
		id:          id,
		conn:        conn,
		proxyBucket: common.NewBucket(200),
	}
}

func (sc *SClient) Serve() {
	for {
		if data, err := DecodeReceive(sc.conn); err == nil {
			switch data[0] {
			case 1:
				switch data[1] {
				case 2:
					if ps := sc.proxyBucket.Get(common.XidToString(data[2:14])); ps != nil {
						if rc := ps.(*ProxyServer).ResultBucket.Get(common.XidToString(data[14:26])); rc != nil {
							if data[26] == 1 {
								rc.(chan bool) <- true
							} else {
								rc.(chan bool) <- false
							}
						} else {
							log.Warn("cant find recive creat tr result channel")
						}
					} else {
						log.Warn("new transport rep not found ps")
					}
				}
			case 3:
				if ps := sc.proxyBucket.Get(common.XidToString(data[1:13])); ps != nil {
					if tr := ps.(*ProxyServer).TransportBucket.Get(common.XidToString(data[13:25])); tr != nil {
						tr.(*Transport).Write(data[25:])
					} else {
						log.Warn("cant find tr exit...")
					}
				} else {
					log.Warn("cant find ps exit..")
				}
			case 4:
				if ps, err := sc.AddProxy(common.AddrByteToString(data[1:]), "temp", "", true); err != nil {
					EncodeSend(sc.conn, []byte{1, 3, 0})
					log.Warn("client request create temp proxy failed..: ", err)
				} else {
					lp := make([]byte, 2)
					binary.BigEndian.PutUint16(lp, ps.ListenPort)
					if _, err := EncodeSend(sc.conn, append([]byte{1, 3, 1}, lp...)); err != nil {
						log.Warn("send temp proxy result err ", err)
					}
				}
			case 5:
				if ps := sc.proxyBucket.Get(common.XidToString(data[1:13])); ps != nil {
					if tr := ps.(*ProxyServer).TransportBucket.Get(common.XidToString(data[13:])); tr != nil {
						tr.(*Transport).Close(false)
						ps.(*ProxyServer).TransportBucket.Remove(common.XidToString(data[13:]))
					} else {
						log.Warn("close tr not found")
					}
				} else {
					log.Warn("close pr not found")
				}
			default:
				log.Warn("not supported cmd faild ", data[0])
				return
			}
		} else {
			log.Warn("receive client data faild , close client ")
			return
		}
	}
}

func (sc *SClient) AddProxy(dest, mark, listenPort string, isTemp bool) (*ProxyServer, error) {
	if destAddr, err := common.AddrStringToByte(dest, "tcp"); err != nil {
		return nil, err
	} else {
		pid := xid.New()
		if ps, err := NewProxyServer(pid.Bytes(), mark, isTemp, sc.conn, destAddr, listenPort); err != nil {
			return nil, err
		} else {
			sc.proxyBucket.Set(pid.String(), ps)
			go ps.Serve()
			return ps, nil
		}
	}
}

func (sc *SClient) DelProxy(pid string) error {
	if ps := sc.proxyBucket.Get(pid); ps != nil {
		ps.(*ProxyServer).Close()
		return nil
	}
	return errors.New("cant find proxyServer by id: " + pid)
}
