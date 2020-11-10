package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gitlab.com/inview-team/raptor_team/worker/internal/rabbit"
)

func main() {
	logger.Info.Print("Worker start")
	rabbit.RabbitRun()
}
