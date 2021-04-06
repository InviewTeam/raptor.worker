package main

import (
	"context"
	"gitlab.com/inview-team/raptor_team/worker/internal/worker"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
)

var (
	rabbitAddr  = os.Getenv("RABBIT_ADDR")
	rabbitQueue = os.Getenv("RABBIT_QUEUE")
)

func main() {
	logger.Info.Print("Worker start")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker := worker.New(rabbitAddr, rabbitQueue)
	go worker.Run(ctx)

	for {
		select {
		case <-done:
			signal.Stop(done)
			return
		}
	}
}
