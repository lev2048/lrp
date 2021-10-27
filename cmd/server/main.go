package main

import (
	"flag"
	"lrp/internal/lrp"
	"os"
	"os/signal"
	"syscall"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log        *logrus.Logger
	token      int
	proto      string
	serverAddr string
	webSerAddr string
)

func init() {
	log = logrus.New()
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
	log.SetOutput(os.Stdout)

	flag.IntVar(&token, "token", 0, "authentication key")
	flag.StringVar(&proto, "proto", "tcp", "protocol [ tcp , kcp , quic ]")
	flag.StringVar(&serverAddr, "server", "0.0.0.0:8801", "input server addr")
	flag.StringVar(&webSerAddr, "web", "0.0.0.0:8802", "input web server addr")
}

func main() {
	flag.Parse()
	server := lrp.NewServer()
	if tk, ok := server.Run(serverAddr, proto, webSerAddr, uint32(token)); !ok {
		log.Error("server start failed exit")
		return
	} else {
		log.Info("server start successful")
		log.Info("access token: ", tk)
		log.Info("server proto: ", proto)
		log.Info("server  addr: ", serverAddr)
		log.Info("webser  addr: ", webSerAddr)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	log.Info("server shutdown ...")
	log.Info("exit")
}
