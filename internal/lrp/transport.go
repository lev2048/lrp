package lrp

import (
	nt "lrp/internal/conn"
	"net"
)

type Transport struct {
	id   []byte
	cc   nt.Conn
	conn net.Conn
	exit chan bool
}

func NewTransport(id []byte, cc nt.Conn, conn net.Conn) *Transport {
	return &Transport{
		id:   id,
		cc:   cc,
		conn: conn,
		exit: make(chan bool),
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
				payload := append([]byte{3}, tr.id...)
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
