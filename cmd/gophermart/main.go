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

	var wg sync.WaitGroup
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	go func() {
		wg.Add(1)
		if err := app.Consume(ctx, cfg); err != nil {
			log.Fatal(err)
		}
		defer wg.Done()
	}()

	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
	wg.Wait()
}
