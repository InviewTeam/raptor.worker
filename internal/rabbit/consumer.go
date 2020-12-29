package rabbit

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gitlab.com/inview-team/raptor_team/worker/internal/structures/task"
	"os"
)

var (
	rabbitLogin   = os.Getenv("RABBIT_LOGIN")
	rabbitPwd     = os.Getenv("RABBIT_PASSWORD")
	rabbitHost    = os.Getenv("RABBIT")
	rabbitPort    = os.Getenv("RABBIT_PORT")
	rabbitChannel = os.Getenv("RABBIT_CHANNEL")
)

func failOnError(err error, msg string) {
	if err != nil {
		logger.Critical.Fatalf("%s:%s", msg, err)
	}
}

func RabbitRun() {
	rabbitAddr := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitLogin, rabbitPwd, rabbitHost, rabbitPort)
	conn, err := amqp.Dial(rabbitAddr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitChannel,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			logger.Info.Printf(string(msg.Body))

			taskInfo := &task.Task{}
			err := json.Unmarshal(msg.Body, taskInfo)

			if err != nil {
				logger.Error.Printf(err.Error())
			}

			if taskInfo.UUID != "" && taskInfo.CameraIP != "" && taskInfo.Status == "" {
				logger.Info.Printf("Start new task %s", taskInfo.UUID)

			} else if taskInfo.UUID != "" && taskInfo.Status != "" && taskInfo.CameraIP == "" {
				logger.Info.Printf("Stop task %s", taskInfo.UUID)
			} else {
				logger.Error.Printf("Unsupported format")
			}

		}
	}()

	logger.Info.Printf(" [*] Waiting for messages.")
	<-forever
}
