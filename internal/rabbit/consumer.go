package rabbit

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/inview-team/raptor_team/worker/internal/cameras"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gitlab.com/inview-team/raptor_team/worker/internal/structures"
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
	rabbitAddr := fmt.Sprintf("amqps://%s:%s@%s", rabbitLogin, rabbitPwd, rabbitHost)
	isStream := false
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
	stream := make(chan bool)

	go func() {
		for msg := range msgs {
			logger.Info.Printf(string(msg.Body))

			taskInfo := &structures.Task{}
			err := json.Unmarshal(msg.Body, taskInfo)

			if err != nil {
				logger.Error.Printf(err.Error())
			}

			if taskInfo.Status == "" {
				if !isStream {
					logger.Info.Printf("Start new task %s", taskInfo.UUID)
					isStream = true
					go cameras.WorkWithVideo(taskInfo.CameraIP, taskInfo.ADDR, stream)
				} else {
					logger.Info.Println("Stream already started")
				}
			} else if taskInfo.CameraIP == "" {
				if isStream {
					logger.Info.Printf("Stop task %s", taskInfo.UUID)
					isStream = false
					stream <- true
				} else {
					logger.Info.Println("Stream didn't started")
				}
			} else {
				logger.Error.Printf("Unsupported format")
			}
		}
	}()

	logger.Info.Printf(" [*] Waiting for messages.")
	<-forever
}
