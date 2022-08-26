package main

import (
	"context"
	"getblock/internal/server"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, finish := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func(ctx context.Context) {
		s := server.New(server.NewCofig())
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	<-sigCh
	finish()
}
