package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

var (
	url = "rtsp://user:qwerty1234@10.10.0.136:5504/cam/realmonitor?channel=1&subtype=0"
)

func main() {
	logger.Info.Print("Worker start")

	cameras.WorkWithVideo(url)
	cameras.NewPeerConnection()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Info.Println(sig)
		done <- true
	}()
	logger.Info.Println("Server start Awaiting Signal")
	<-done
	logger.Info.Println("Worker stopped")
}
