package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"os"
)

var (
	rabbitAddr  = os.Getenv("RABBIT_ADDR")
	rabbitQueue = os.Getenv("RABBIT_QUEUE")
)

func main() {
	logger.Info.Println("Start worker service")

	/* done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker := worker.New(rabbitAddr, rabbitQueue)

	go worker.Listen(ctx)
	go worker.Run(ctx)

	for {
		select {
		case <-done:
			signal.Stop(done)
			return
		}
	} */
	test := make(chan struct{})
	cameras.GetImagesFromRTSP("rtsp://user:qwerty1234@10.10.0.136:5504/cam/realmonitor?channel=1&subtype=0", "test", test)
}
