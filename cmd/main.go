package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gitlab.com/inview-team/raptor_team/worker/internal/worker"
)

var (
	rabbit_addr  = os.Getenv("RABBIT_ADDR")
	rabbit_queue = os.Getenv("RABBIT_QUEUE")
)

func main() {
	logger.Info.Print("Worker start")
	http.HandleFunc("/doSignaling", cameras.DoSignaling)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker := worker.New(rabbit_addr, rabbit_queue)
	go worker.Run(ctx)

	for {
		select {
		case <-done:
			signal.Stop(done)
			return
		}
	}
}
