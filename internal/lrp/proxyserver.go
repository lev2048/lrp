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
	id              []byte
	exit            chan bool
	isOk            chan bool
	Conn            nt.Conn
	DestAddr        []byte
	Listener        net.Listener
	ListenPort      uint16
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
		id:              id,
		exit:            make(chan bool),
		isOk:            make(chan bool),
		Conn:            conn,
		DestAddr:        dest,
		Listener:        ln,
		ListenPort:      uint16(lp),
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
	tid := xid.New().Bytes()
	data := append([]byte{2}, tid...)
	data = append(data, ps.DestAddr...)
	if _, err := EncodeSend(ps.Conn, data); err != nil {
		log.Warn("send accept pk to client error", err)
		return
	}
	defer close(ps.isOk)
	for {
		select {
		case <-time.After(time.Second * 15):
			log.Warn("wait accept pk reply timeout")
			conn.Close()
			return
		case isRemoteOk := <-ps.isOk:
			if isRemoteOk {
				tr := NewTransport(tid, ps.Conn, conn)
				ps.TransportBucket.Set(XidToString(tid), tr)
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
		v.(*Transport).Close()
	}
	ps.TransportBucket = nil
	close(ps.exit)
	//todo: remove stat
}
