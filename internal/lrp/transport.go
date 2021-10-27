package lrp

import (
	nt "lrp/internal/conn"
	"net"
)

type Transport struct {
	id       []byte
	pd       []byte
	cc       nt.Conn
	conn     net.Conn
	exit     chan bool
	isServer bool
}

func NewTransport(isServer bool, id, pd []byte, cc nt.Conn, conn net.Conn) *Transport {
	return &Transport{
		id:       id,
		pd:       pd,
		cc:       cc,
		conn:     conn,
		exit:     make(chan bool),
		isServer: isServer,
	}
}

func (tr *Transport) Serve() {
	defer tr.conn.Close()
	buf := make([]byte, 1460)
	for {
		select {
		case <-tr.exit:
			log.Info("transport has exited")
			return
		default:
			if n, err := tr.conn.Read(buf); err != nil {
				//todo: notify client close tr
				log.Warn("recive data err", err)
				return
			} else {
				var payload []byte
				if tr.isServer {
					payload = append([]byte{3}, tr.id...)
				} else {
					payload = append([]byte{3}, tr.pd...)
					payload = append(payload, tr.id...)
				}
				payload = append(payload, buf[:n]...)
				if _, err := EncodeSend(tr.cc, payload); err != nil {
					log.Warn("send data to client err", err)
					return
				}
			}
		}
	}
}

func (tr *Transport) Write(data []byte) error {
	if _, err := tr.conn.Write(data); err != nil {
		log.Warn("send data to user err", err)
		tr.Close()
		return err
	}
	return nil
}

func (tr *Transport) Close() {
	close(tr.exit)
}
