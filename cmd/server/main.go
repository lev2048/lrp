package main

import (
	"flag"
	"lrp/internal/api"
	"lrp/internal/lrp"
	"lrp/internal/status"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	Secret          string = "gms90k3g6ej3"
	AppMode         string = "release"
	ApiReadTimeout  int    = 120
	ApiWriteTimeout int    = 120
)

var (
	log        *logrus.Logger
	sig        chan os.Signal = make(chan os.Signal, 1)
	token      int
	proto      string
	monitor    *status.Monitor = status.NewMonitor()
	serverPort int
	webSerPort int
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
	flag.IntVar(&serverPort, "server", 8801, "input server addr")
	flag.IntVar(&webSerPort, "web", 0, "input web server addr")
}

func main() {
	flag.Parse()
	server := lrp.NewServer()
	if tk, ok := server.Run(strconv.Itoa(serverPort), proto, uint32(token)); !ok {
		log.Error("server start failed exit")
		return
	} else {
		log.Info("server start successful")
		log.Info("access token: ", tk)
		log.Info("server proto: ", proto)
		log.Info("server  port: ", serverPort)
	}
	if webSerPort != 0 {
		sp := strconv.Itoa(webSerPort)
		log.Info("web server start ...")
		gin.SetMode(AppMode)
		engine := gin.New()
		api := api.NewApi(server, engine, monitor, Secret)
		httpServer := &http.Server{
			Addr:         ":" + sp,
			Handler:      engine,
			ReadTimeout:  time.Duration(ApiReadTimeout) * time.Second,
			WriteTimeout: time.Duration(ApiWriteTimeout) * time.Second,
		}
		api.SetRouter()
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Warn("web server start faild: ", err)
				sig <- syscall.SIGTERM
			}
		}()
		monitor.Start()
		log.Info("web server start successful ")
		log.Info("web url: ", "http://"+server.ExternalIp+":"+sp)
	}
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	log.Info("server shutdown ...")
	//todo close server
	if webSerPort != 0 {
		log.Info("close monitor ...")
		monitor.Stop()
	}
	log.Info("done")
}
