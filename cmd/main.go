package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gitlab.com/inview-team/raptor_team/worker/internal/rabbit"
	"net/http"
)

func main() {
	logger.Info.Print("Worker start")
	http.HandleFunc("/doSignaling", cameras.DoSignaling)
	rabbit.RabbitRun()
}
