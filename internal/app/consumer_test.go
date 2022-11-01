package app

import (
	"context"
	"github.com/syols/go-devops/config"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestConsumeHandler(t *testing.T) {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	cfg := config.Config{}
	errs := make(chan error, 1)
	Consume(ctx, cfg, errs)
}

