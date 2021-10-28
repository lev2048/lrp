package lrp

import (
	"errors"
	"io"
	nt "lrp/internal/conn"
	"net"
)

type Transport struct {
	id       []byte
	pd       []byte
	cc       nt.Conn
	conn     net.Conn
	exit     chan bool
	isClosed bool
	isServer bool
}

func NewTransport(isServer bool, id, pd []byte, cc nt.Conn, conn net.Conn) *Transport {
	return &Transport{
		id:       id,
		pd:       pd,
		cc:       cc,
		conn:     conn,
		exit:     make(chan bool),
		isClosed: false,
		isServer: isServer,
	}
}

func (tr *Transport) Serve() {
	buf := make([]byte, 1460)
	for {
		select {
		case <-tr.exit:
			return
		default:
			if n, err := tr.conn.Read(buf); err != nil {
				if tr.isClosed {
					return
				}
				if err != io.EOF {
					log.Warn("recive data err ", err)
				}
				tr.Close(true)
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
					log.Warn("send data to client err ", err)
					tr.Close(false)
				}
			}
		}
	}
}

func (tr *Transport) Write(data []byte) error {
	if tr.isClosed {
		return errors.New("transport is closed")
	}
	if _, err := tr.conn.Write(data); err != nil {
		log.Warn("send data to user err", err)
		tr.Close(false)
		return err
	}
	return nil
}

func (tr *Transport) Close(notify bool) {
	if tr.isClosed {
		return
	}
	if notify {
		pl := []byte{5}
		if tr.isServer {
			pl = append(pl, tr.id...)
		} else {
			pl = append(append(pl, tr.pd...), tr.id...)
		}
		if _, err := EncodeSend(tr.cc, pl); err != nil {
			log.Warn("send close data err ", err)
		}
	}
	close(tr.exit)
	tr.isClosed = true
	tr.conn.Close()
}
