package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis"
	"github.com/microlib/simple"
)

var (
	logger simple.Logger
	client *redis.Client
	config Config
)

func startHttpServer(port string) *http.Server {
	srv := &http.Server{Addr: ":" + port}

	// add all the routes and link to handlers
	http.HandleFunc("/iotdata", IotcontrollerPostIOTData)
	http.HandleFunc("/iotdata/list", ListIOTData)
	http.HandleFunc("/isalive", IsAlive)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("Httpserver: ListenAndServe() error: %s", err.Error()))
		}
	}()

	return srv
}

func cleanup() {
	logger.Info("main: server cleanup")
}

func main() {
	var cfg Config
	config = cfg.Init("config.json")
	logger.Level = config.Level

	// checkif we have caching enabled
	if config.Cache == "true" {
		server, err := config.GetServer("redis")
		if err != nil {
			logger.Error(err.Error())
		}

		client = redis.NewClient(&redis.Options{
			Addr:     server.Host + ":" + server.Port,
			Password: "",
			DB:       0,
		})

		logger.Info(client.Ping().String())
		logger.Info(fmt.Sprintf("Redis info %s %s", server.Host, server.Port))
	}

	srv := startHttpServer(config.Port)
	logger.Info(fmt.Sprintf("main: starting server on port %s", srv.Addr))
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan
	cleanup()
	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("main: server shutdown successfully")
	os.Exit(code)
}
