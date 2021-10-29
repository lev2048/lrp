package lrp

import (
	"encoding/binary"
	"io"
	"lrp/internal/conn"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

const (
	version string = "v1.0.0"
)

func init() {
	log = logrus.New()
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
	log.SetOutput(os.Stdout)
}

type ProxyInfo struct {
	Id      string `json:"id"`
	Info    string `json:"info"`
	Mark    string `json:"mark"`
	Status  int    `json:"status"`
	ConnNum int    `json:"connNum"`
}

type ClientInfo struct {
	Id         string      `json:"id"`
	Mark       string      `json:"mark"`
	Online     bool        `json:"online"`
	ProxyInfos []ProxyInfo `json:"proxyInfos"`
}

type ServerInfo struct {
	Version     string       `json:"version"`
	Protocol    string       `json:"protocol"`
	ProxyNum    int          `json:"proxyNum"`
	TProxyNum   int          `json:"tProxyNum"`
	ConnTotal   int          `json:"connTotal"`
	ClientNum   int          `json:"clientNum"`
	ExternalIp  string       `json:"externalIp"`
	ServerPort  string       `json:"serverPort"`
	ClientInfos []ClientInfo `json:"clientInfos"`
}

func EncodeSend(c conn.Conn, data []byte) (int, error) {
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	copy(buf[4:], data)
	n, err := c.Write(buf)
	return n, err
}

func DecodeReceive(r io.Reader) ([]byte, error) {
	len := make([]byte, 4)
	if _, err := io.ReadFull(r, len); err != nil {
		return nil, err
	}
	payload := make([]byte, binary.BigEndian.Uint32(len))
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
