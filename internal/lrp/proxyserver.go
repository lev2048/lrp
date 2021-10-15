package lrp

import (
	"lrp/internal/common"
	nt "lrp/internal/conn"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
)

type ProxyServer struct {
	Id              []byte
	Ch              chan bool
	Conn            nt.Conn
	DestAddr        []byte
	Listener        net.Listener
	ListenPort      uint32
	TransportBucket *common.Bucket
}

func NewProxyServer(id []byte, conn nt.Conn, dest []byte) (*ProxyServer, error) {
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}
	lp, err := strconv.Atoi(strings.Split(ln.Addr().String(), ":")[1])
	if err != nil {
		return nil, err
	}
	return &ProxyServer{
		Id:              id,
		Ch:              make(chan bool, 0),
		Conn:            conn,
		Listener:        ln,
		ListenPort:      uint32(lp),
		TransportBucket: common.NewBucket(1024),
	}, nil
}

func (ps *ProxyServer) Serve() {
	defer ps.Listener.Close()
	for {
		if conn, err := ps.Listener.Accept(); err != nil {
			log.Warn("ProxyServer accept conn faild", err)
			return
		} else {
			go ps.handleConn(conn)
		}
	}
}

func (ps *ProxyServer) handleConn(conn net.Conn) {
	tid := xid.New()
	data := append([]byte{2}, tid.Bytes()...)
	data = append(data, ps.DestAddr...)
	if _, err := EncodeSend(ps.Conn, data); err != nil {
		log.Warn("send accept pk to client error", err)
		return
	}
	for {
		select {
		case <-time.After(time.Second * 15):
			log.Warn("wait accept pk reply timeout")
			return
		case isRemoteOk := <-ps.Ch:
			if isRemoteOk {
				tr := NewTransport()
				ps.TransportBucket.Set(string(tid.Bytes()), tr)
			} else {
				log.Warn("dest obj connect failed")
				conn.Close()
			}
			return
		}
	}
}

func (ps *ProxyServer) Close() {

}
