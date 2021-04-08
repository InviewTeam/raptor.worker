package cameras

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"gitlab.com/inview-team/raptor_team/worker/internal/kafka"
	"image"
	"os"
	"time"

	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gocv.io/x/gocv"
)

type Result struct {
	Pix      []byte `json:"pix"`
	Channels int    `json:"channels"`
	Rows     int    `json:"rows"`
	Cols     int    `json:"cols"`
}

func GetImagesFromRTSP(cameraUrl string, topic string) {
	camera, err := gocv.OpenVideoCapture(cameraUrl)
	if err != nil {
		logger.Error.Panic("Error in opening camera: " + err.Error())
	}

	var brokers = []string{os.Getenv("KAFKAPORT")}
	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		panic("Failed to connect to Kafka. Error: " + err.Error())
	}
	//Close producer to flush(i.e., push) all batched messages into Kafka queue
	defer func() { producer.Close() }()

	frame := gocv.NewMat()
	for {
		if !camera.Read(&frame) {
			continue
		}

		imgInterface, err := frame.ToImage()
		if err != nil {
			logger.Error.Panic(err.Error())
		}

		img, ok := imgInterface.(*image.RGBA)
		if !ok {
			logger.Error.Panic("Type assertion of pic (type image.Image interface) to type image.RGBA failed")
		}

		doc := Result{
			Pix:      img.Pix,
			Channels: frame.Channels(),
			Rows:     frame.Rows(),
			Cols:     frame.Cols(),
		}

		docBytes, err := json.Marshal(doc)
		if err != nil {
			logger.Critical.Fatal("Json marshalling error. Error:", err.Error())
		}

		msg := &sarama.ProducerMessage{
			Topic:     topic,
			Value:     sarama.ByteEncoder(docBytes),
			Timestamp: time.Now(),
		}

		producer.Input() <- msg
	}
}
