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
	"unsafe"
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

//AddrStringToByte 网络地址 “1.1.1.1:80” -> {1,1,1,1,11,11} 端口大端序
func AddrStringToByte(addr string, netType string) ([]byte, error) {
	if netType == "tcp" {
		address, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			fmt.Println("地址格式转换失败")
			return nil, err
		}
		port := make([]byte, 2)
		binary.BigEndian.PutUint16(port, uint16(address.Port))
		result := append(address.IP.To4(), port...)
		return result, nil
	} else {
		address, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Println("地址格式转换失败")
			return nil, err
		}
		port := make([]byte, 2)
		binary.BigEndian.PutUint16(port, uint16(address.Port))
		result := append(address.IP.To4(), port...)
		return result, nil
	}
}

func XidToString(id []byte) string {
	text := make([]byte, 20)
	encode(text, id[:])
	return *(*string)(unsafe.Pointer(&text))
}

func encode(dst, id []byte) {
	encoding := "0123456789abcdefghijklmnopqrstuv"
	_ = dst[19]
	_ = id[11]

	dst[19] = encoding[(id[11]<<4)&0x1F]
	dst[18] = encoding[(id[11]>>1)&0x1F]
	dst[17] = encoding[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[16] = encoding[id[10]>>3]
	dst[15] = encoding[id[9]&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[11] = encoding[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[9] = encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[7] = encoding[id[4]&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[3] = encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[1] = encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[0] = encoding[id[0]>>3]
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
