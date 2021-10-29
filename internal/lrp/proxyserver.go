package lrp

import (
	"io"
	"lrp/internal/common"
	nt "lrp/internal/conn"
	"net"
	"time"

	"github.com/rs/xid"
)

type ProxyServer struct {
	Id              []byte
	exit            chan bool
	Conn            nt.Conn
	Mark            string
	Temp            bool
	Status          int
	DestAddr        []byte
	IsClosed        bool
	Listener        net.Listener
	ListenPort      uint16
	ResultBucket    *common.Bucket
	TransportBucket *common.Bucket
}

func NewProxyServer(id []byte, mark string, isTemp bool, conn nt.Conn, dest []byte, listenPort string) (*ProxyServer, error) {
	addr := "0.0.0.0:"
	if listenPort != "" {
		addr += listenPort
	} else {
		addr += "0"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &ProxyServer{
		Id:              id,
		exit:            make(chan bool),
		Conn:            conn,
		Mark:            mark,
		Temp:            isTemp,
		Status:          1,
		DestAddr:        dest,
		IsClosed:        false,
		Listener:        ln,
		ListenPort:      uint16(ln.Addr().(*net.TCPAddr).Port),
		ResultBucket:    common.NewBucket(1024),
		TransportBucket: common.NewBucket(1024),
	}, nil
}

func (ps *ProxyServer) Serve() {
	for {
		select {
		case <-ps.exit:
			return
		default:
			if conn, err := ps.Listener.Accept(); err != nil {
				if ps.IsClosed {
					return
				}
				if err != io.EOF {
					log.Warn("ProxyServer accept conn faild", err)
				}
				return
			} else {
				go ps.handleConn(conn)
			}
		}
	}
}

func (ps *ProxyServer) handleConn(conn net.Conn) {
	tid, seq := xid.New().Bytes(), xid.New()
	data := append(append([]byte{2}, ps.Id...), tid...)
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
	ps.TransportBucket, ps.IsClosed = nil, true
	ps.Listener.Close()
	close(ps.exit)
}
