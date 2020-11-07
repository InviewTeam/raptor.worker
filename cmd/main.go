package main

import (
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gitlab.com/inview-team/raptor_team/worker/internal/registry"
)

func main() {
	logger.Info.Print("Worker start")
	task, err := registry.GetNewTask()
	if err != nil {
		logger.Error.Panic(err)
	}
	cameraIP := task.CameraIP
	uuid := task.UUID
	jobs := task.Jobs

	for {
		result := registry.GetStatusOfTask(uuid)
		if !result {
			break
		}
	}

}
