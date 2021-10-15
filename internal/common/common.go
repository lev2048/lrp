package common

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"strconv"
)

func GenerateServerToken() string {
	token := bytes.NewBufferString("")
	for i := 0; i < 6; i++ {
		v, _ := rand.Int(rand.Reader, big.NewInt(9))
		token.WriteString(strconv.Itoa(int(v.Int64())))
	}
	return token.String()
}

//AddrStringToByte 网络地址 {1,1,1,1,11,11} -> “1.1.1.1:80” 端口大端序
func AddrByteToString(addr []byte) string {
	ip := net.IP(addr[0:4])
	port := binary.BigEndian.Uint16(addr[4:6])
	destAddr := ip.String() + ":" + fmt.Sprint(port)
	return destAddr
}

func GenerateTLSConfig(title string) (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{title},
	}, nil
}
