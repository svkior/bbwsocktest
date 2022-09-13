package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/svkior/bbwsocktest/internal/config"
	"github.com/svkior/bbwsocktest/internal/wsclient"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewBBConfig()

	bbClient := wsclient.NewWSClient(cfg)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := bbClient.Init(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	// Starting the waiting goroutine for terminate goroutine
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	g.Go(func() error {
		<-c
		log.Printf("Stoping service because of SIGTERM")
		cancel()
		return nil
	})

	err := g.Wait()
	if err != nil {
		log.Printf("Error running some services %s", err.Error())
	}

	log.Printf("Close")
}
