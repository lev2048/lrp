package lrp

import (
	"lrp/internal/common"
	nt "lrp/internal/conn"
	"net"
	"time"

	"github.com/rs/xid"
)

type ProxyServer struct {
	id              []byte
	exit            chan bool
	Conn            nt.Conn
	DestAddr        []byte
	Listener        net.Listener
	ListenPort      uint16
	ResultBucket    *common.Bucket
	TransportBucket *common.Bucket
}

func NewProxyServer(id []byte, conn nt.Conn, dest []byte) (*ProxyServer, error) {
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}
	return &ProxyServer{
		id:              id,
		exit:            make(chan bool),
		Conn:            conn,
		DestAddr:        dest,
		Listener:        ln,
		ListenPort:      uint16(ln.Addr().(*net.TCPAddr).Port),
		ResultBucket:    common.NewBucket(1024),
		TransportBucket: common.NewBucket(1024),
	}, nil
}

func (ps *ProxyServer) Serve() {
	defer ps.Listener.Close()
	for {
		select {
		case <-ps.exit:
			log.Info("proxy is closed")
			return
		default:
			if conn, err := ps.Listener.Accept(); err != nil {
				log.Warn("ProxyServer accept conn faild", err)
				return
			} else {
				go ps.handleConn(conn)
			}
		}
	}
}

func (ps *ProxyServer) handleConn(conn net.Conn) {
	tid, seq := xid.New().Bytes(), xid.New()
	data := append(append([]byte{2}, ps.id...), tid...)
	data = append(append(data, seq.Bytes()...), ps.DestAddr...)
	if _, err := EncodeSend(ps.Conn, data); err != nil {
		log.Warn("send accept pk to client error", err)
		return
	}
	result := make(chan bool)
	ps.ResultBucket.Set(seq.String(), result)
	for {
		select {
		case <-time.After(time.Second * 15):
			log.Warn("wait accept pk reply timeout")
			conn.Close()
			return
		case isRemoteOk := <-result:
			defer close(result)
			defer ps.ResultBucket.Remove(seq.String())
			if isRemoteOk {
				tr := NewTransport(true, tid, nil, ps.Conn, conn)
				ps.TransportBucket.Set(common.XidToString(tid), tr)
				go tr.Serve()
			} else {
				log.Warn("dest obj connect failed")
				conn.Close()
			}
			return
		}
	}
}

func (ps *ProxyServer) Close() {
	for _, v := range ps.TransportBucket.GetAll() {
		v.(*Transport).Close(true)
	}
	ps.TransportBucket = nil
	close(ps.exit)
}
