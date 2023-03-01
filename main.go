package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/shivendra195/supplyChainManagement/server"

	"github.com/sirupsen/logrus"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.SrvInit()
	go srv.Start()

	<-done
	logrus.Info("Graceful shutdown")
	srv.Stop()
}
