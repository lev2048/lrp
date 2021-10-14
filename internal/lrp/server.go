package lrp

import (
	"errors"
	"lrp/internal/common"
	"lrp/internal/conn"
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
		return
	}

	sc := newSClient(conn)
	go sc.Serve()

	id := xid.New()
	s.clients.Set(id.String(), sc)
}

func (s *Server) StartWebServer() error {
	return nil
}

func (s *Server) AddProxy(cid, dest string) error {
	return nil
}

func (s *Server) DelProxy(cid, pid string) error {
	return nil
}

func (s *Server) GetClientList() {}

type SClient struct {
	conn        conn.Conn
	proxyBucket *common.Bucket
}

func newSClient(conn conn.Conn) *SClient {
	return &SClient{
		conn:        conn,
		proxyBucket: common.NewBucket(1024),
	}
}

func (sc *SClient) Serve() {
	for {
		if data, err := DecodeReceive(sc.conn); err != nil {

		}
	}
}

func (sc *SClient) AddProxy(dest string) error {
	return nil
}

func (sc *SClient) DelProxy(pid string) error {
	return nil
}
