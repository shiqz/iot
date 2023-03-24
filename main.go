package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	_ "iot/internal/pkgs/logger"
	"net"
	"os"
	"os/signal"
	"time"
)

type CollectorServer struct {
	listen net.Listener
}

func (c *CollectorServer) Run() {
	log.Infof("collector running on:%s", c.listen.Addr().String())
	for {
		conn, err := c.listen.Accept()
		if err != nil {
			return
		}
		go func() {
			log.WithFields(log.Fields{"conn": conn.RemoteAddr().String()}).Info("connected")
			defer func() {
				log.WithFields(log.Fields{"conn": conn.RemoteAddr().String()}).Info("disconnected")
			}()
			for {
				var buff = make([]byte, 1024)
				_, err = conn.Read(buff)
				if err != nil {
					if _, ok := err.(net.Error); ok || err == io.EOF {
						log.Errorln(fmt.Sprintf("client net err:%v", err))
						return
					}
					log.Warning(fmt.Sprintf("解析数据包异常 => %v", err))
					continue
				}
				log.Infof("message=>%s", string(buff))
			}
		}()
	}
}

func (c *CollectorServer) Shutdown() error {
	return c.listen.Close()
}

func RunCollectorServer() (srv *CollectorServer, err error) {
	addr := fmt.Sprintf("0.0.0.0:8000")
	srv = &CollectorServer{}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	srv.listen = listen
	return
}

func main() {
	srv, err := RunCollectorServer()
	if err != nil {
		panic(err)
	}
	go srv.Run()

	// 程序退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
