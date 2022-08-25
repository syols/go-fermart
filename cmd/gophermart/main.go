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
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 1)
	wg.Add(1)
	go app.Consume(ctx, cfg, errs)

	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	server.Run()
	wg.Wait()
	
	for err := range errs { // TODO
		log.Print(err)
	}
}
