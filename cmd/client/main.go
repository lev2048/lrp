package main

import (
	"flag"
	"lrp/internal/lrp"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log        *logrus.Logger
	token      int
	proto      string
	verbose    bool
	destAddr   string
	serverAddr string
)

func init() {
	log = logrus.New()
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
	log.SetOutput(os.Stdout)

	flag.IntVar(&token, "token", 0, "authentication key")
	flag.BoolVar(&verbose, "vbose", false, "verbose mode")
	flag.StringVar(&proto, "proto", "tcp", "protocol [ tcp , kcp , quic ]")
	flag.StringVar(&serverAddr, "server", "", "input server addr")
	flag.StringVar(&destAddr, "dest", "", "Enter the target address you want to proxy eg: 192.168.1.1:80")
}

func main() {
	flag.Parse()
	if token == 0 || serverAddr == "" {
		log.Error("token or serverAddress is required")
		return
	}
	client := lrp.NewClient(uint32(token), verbose)
	log.Info("connect server ...")
	if err := client.Run(serverAddr, proto); err != nil {
		log.Error("connect failed: ", err)
		return
	}
	if destAddr != "" {
		log.Info("request proxy port ...")
		if port, err := client.AddTempProxy(destAddr); err != nil {
			log.Error("request failed ", err)
			return
		} else {
			log.Info("request successful")
			log.Info("server: " + strings.Split(serverAddr, ":")[0] + ":" + strconv.Itoa(int(port)) + " => " + destAddr)
		}
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	log.Info("client shutdown ...")
	log.Info("exit")
}
