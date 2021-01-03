package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Info.Print("Worker start")
	var url = "rtsp://user:qwerty1234@10.10.0.136:5506/cam/realmonitor?channel=1&subtype=0"
	go cameras.ServeStream(":8083")
	go cameras.SubsribeStream(url)
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()
	log.Println("Server Start Awaiting Signal")
	<-done
	log.Println("Exiting")
}
