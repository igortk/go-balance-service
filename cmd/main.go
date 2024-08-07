package main

import (
	"balance-service/config"
	"balance-service/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	svr, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	svr.Init()
	svr.Run()
}
