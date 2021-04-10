package worker

import (
	"context"
	"encoding/json"
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"log"

	"gitlab.com/inview-team/raptor_team/worker/internal/rabbit"
	"gitlab.com/inview-team/raptor_team/worker/internal/structures"
)

type Worker struct {
	con           *rabbit.Consumer
	tasksIncoming chan []byte
	tasksInWork   map[string]chan struct{}
}

func New(addr, queue string) *Worker {
	tasks := make(chan []byte)
	return &Worker{
		con:           rabbit.NewConsumer(addr, queue, tasks),
		tasksIncoming: tasks,
	}
}

func (w *Worker) Listen(ctx context.Context) error {
	log.Printf("start listening")
	return w.con.Receive(ctx)
}

func (w *Worker) Run(ctx context.Context) error {
	var task structures.Task
	var err error

	for {
		select {
		case data := <-w.tasksIncoming:
			logger.Info.Println(w.tasksIncoming)
			err = json.Unmarshal(data, &task)
			if err != nil {
				log.Printf(err.Error())
				break
			}

			if task.Status == structures.InWork {
				done := make(chan struct{})
				w.tasksInWork[task.UUID] = done
				go cameras.WorkerLoop(task.CameraIP, task.UUID, done)

			} else if task.Status == structures.Stopped {
				close(w.tasksInWork[task.UUID])
				delete(w.tasksInWork, task.UUID)
			}

		case <-ctx.Done():
			return w.Stop()
		}
	}
}

func (w *Worker) Stop() error {
	for _, done := range w.tasksInWork {
		close(done)
	}
	return w.con.Close()
}
