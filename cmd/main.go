package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
)

func main() {
	logger.Info.Print("Worker start")
	var url = "rtsp://user:qwerty1234@10.10.0.136:5506/cam/realmonitor?channel=1&subtype=0"
	cameras.WorkWithVideo(url, ":8080")
}
