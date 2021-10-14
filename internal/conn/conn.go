package conn

import (
	"errors"
	"io"
	"strings"
)

type Conn interface {
	io.ReadWriteCloser
	Name() string
	Info() string
	Dial(dst string) (Conn, error)
	Accept() (Conn, error)
	Listen(dst string) (Conn, error)
}

func NewConn(proto string) (Conn, error) {
	switch strings.ToLower(proto) {
	case "tcp":
		return &TcpConn{}, nil
	case "kcp":
		return &KcpConn{}, nil
	case "quic":
		return &QuicConn{}, nil
	default:
		return nil, errors.New("undefined proto " + proto)
	}
}

func SupportProtos() []string {
	return []string{"tcp", "kcp", "quic"}
}
