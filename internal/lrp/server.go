package lrp

import (
	"encoding/binary"
	"errors"
	"lrp/internal/common"
	nt "lrp/internal/conn"

	"github.com/rs/xid"
)

type Server struct {
	token   string
	clients *common.Bucket
}

func NewServer() *Server {
	return &Server{
		clients: common.NewBucket(1024),
	}
}

func (s *Server) Run(sp, proto, wp, tk string) bool {
	var err error
	defer func() {
		if err != nil {
			log.Warn(err.Error())
		}
	}()

	if tk != "" {
		if tk = common.GenerateServerToken(); tk == "" {
			err = errors.New("generate token failed")
			return false
		}
	}

	ns, err := nt.NewConn(proto)
	if err != nil {
		return false
	}

	ln, err := ns.Listen("0.0.0.0:" + sp)
	if err != nil {
		return false
	}

	go func() {
		for {
			if conn, err := ln.Accept(); err != nil {
				go s.handleClient(conn)
			} else {
				log.Warn("accept new client error:", err)
			}
		}
	}()
	return true
}

func (s *Server) handleClient(conn nt.Conn) {
	if tk, err := DecodeReceive(conn); err != nil || string(tk[1:]) != s.token {
		log.Warn("auth client faild", string(tk[1:]))
		EncodeSend(conn, []byte{1, 0})
		return
	}
	if _, err := EncodeSend(conn, []byte{1, 1}); err != nil {
		log.Warn("send auth reply err", err)
		return
	}
	id, sc := xid.New(), newSClient(conn)
	s.clients.Set(id.String(), sc)
	defer func() {
		conn.Close()
		s.clients.Remove(id.String())
	}()
	sc.Serve()
}

func (s *Server) AddProxy(cid, dest string) error {
	if client := s.clients.Get(cid); client != nil {
		if _, err := client.(*SClient).AddProxy(dest, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) DelProxy(cid, pid string) error {
	if client := s.clients.Get(cid); client != nil {
		if err := client.(*SClient).DelProxy(pid); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) GetClientList() {}

type SClient struct {
	conn        nt.Conn
	proxyBucket *common.Bucket
}

func newSClient(conn nt.Conn) *SClient {
	return &SClient{
		conn:        conn,
		proxyBucket: common.NewBucket(200),
	}
}

func (sc *SClient) Serve() {
	for {
		if data, err := DecodeReceive(sc.conn); err != nil {
			switch data[0] {
			case 3:
				if ps := sc.proxyBucket.Get(string(data[1:9])); ps != nil {
					if tr := ps.(*ProxyServer).TransportBucket.Get(string(data[9:17])); tr != nil {
						tr.(*Transport).Write(data[17:])
					} else {
						log.Warn("cant find tr exit...")
					}
				} else {
					log.Warn("cant find ps exit..")
				}
			case 4:
				if ps, err := sc.AddProxy(common.AddrByteToString(data[1:7]), true); err != nil {
					EncodeSend(sc.conn, []byte{4, 0})
					log.Warn("client request create temp proxy failed..", err)
				} else {
					lp := make([]byte, 2)
					binary.BigEndian.PutUint16(lp, ps.ListenPort)
					if _, err := EncodeSend(sc.conn, append([]byte{4, 1}, lp...)); err != nil {
						log.Warn("send temp proxy result err", err)
					}
				}
			default:
				log.Warn("not supported cmd faild", data[0])
				return
			}
		} else {
			log.Warn("receive client data faild , close client")
			return
		}
	}
}

func (sc *SClient) AddProxy(dest string, isTemp bool) (*ProxyServer, error) {
	if destAddr, err := common.AddrStringToByte("tcp", dest); err != nil {
		return nil, err
	} else {
		pid := xid.New()
		if ps, err := NewProxyServer(pid.Bytes(), sc.conn, destAddr); err != nil {
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
