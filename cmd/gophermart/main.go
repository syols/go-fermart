package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
)

func main() {
	log.SetOutput(os.Stdout)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(cfg)

	var wg sync.WaitGroup
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		err := app.Consume(ctx, &wg, cfg)
		if err != nil {
			log.Fatal(err)
		}
	}()

	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
	wg.Wait()
}
