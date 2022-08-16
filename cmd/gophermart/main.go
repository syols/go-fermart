package main

import (
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
