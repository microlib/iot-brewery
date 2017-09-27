package main

import (
	"crypto/tls"
	"fmt"
	"github.com/robfig/cron"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/microlib/simple"
)

var (
	logger simple.Logger
	config Config
)

func main() {

	var cfg Config
	config = cfg.Init("config.json")
	logger.Level = config.Level

	initSpi()
	time.Sleep(200 * time.Millisecond)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cr := cron.New()
	cr.AddFunc(config.Cron,
		func() {
			processIOT(tr, config)
		})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		cleanup(cr)
		os.Exit(0)
	}()

	cr.Start()
	logger.Info("IOTController started")
	logger.Info(fmt.Sprintf("Golang crontab is : %s ", config.Cron))

	for {
		logger.Trace(fmt.Sprintf("NOP sleeping for %d seconds", config.Sleep))
		time.Sleep(time.Duration(config.Sleep) * time.Second)
	}
}
