package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	url  = os.Getenv("CAMERA_IP")
	addr = os.Getenv("STREAM_PORT")
)

func main() {
	logger.Info.Print("Worker start")

	http.HandleFunc("/doSignaling", cameras.DoSignaling)
	cameras.WorkWithVideo(url, addr)

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
