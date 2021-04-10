package cameras

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/cgo/ffmpeg"
	"gitlab.com/inview-team/raptor_team/worker/internal/kafka"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"log"
	"os"
	"time"

	"github.com/deepch/vdk/format/rtspv2"
)

type Result struct {
	Y  []uint8 `json:"y"`
	Cb []uint8 `json:"cb"`
	Cr []uint8 `json:"cr"`
}

func WorkerLoop(url string, uuid string, done chan struct{}) {
	logger.Info.Println("Start streaming: ", uuid)
	for {
		select {
		case <-done:
			return
		default:
			err := GetImagesFromRTSP(url, uuid, done)
			if err != nil {
				logger.Error.Panic(err.Error())
			}
		}
	}
}

func GetImagesFromRTSP(cameraUrl string, topic string, done chan struct{}) error {
	var brokers = []string{os.Getenv("KAFKAPORT")}
	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		return err
	}

	defer func() { producer.Close() }()

	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: cameraUrl, DisableAudio: true, DialTimeout: 3 * time.Second, ReadWriteTimeout: 3 * time.Second, Debug: false})
	if err != nil {
		return err
	}

	defer RTSPClient.Close()

	var videoIDX int
	for i, codec := range RTSPClient.CodecData {
		if codec.Type().IsVideo() {
			videoIDX = i
		}
	}
	var FrameDecoderSingle *ffmpeg.VideoDecoder

	FrameDecoderSingle, err = ffmpeg.NewVideoDecoder(RTSPClient.CodecData[videoIDX].(av.VideoCodecData))
	if err != nil {
		log.Fatalln("FrameDecoderSingle Error", err)
	}
	for {
		select {
		case <-done:
			break
		default:
			packet := <-RTSPClient.OutgoingPacketQueue
			pic, err := FrameDecoderSingle.DecodeSingle(packet.Data)
			doc := Result{
				Y:  pic.Image.Y,
				Cb: pic.Image.Cb,
				Cr: pic.Image.Cr,
			}

			docBytes, err := json.Marshal(doc)
			if err != nil {
				return err
			}

			msg := &sarama.ProducerMessage{
				Topic:     topic,
				Value:     sarama.ByteEncoder(docBytes),
				Timestamp: time.Now(),
			}

			producer.Input() <- msg
		}
	}
	logger.Info.Println("Stop streaming")
	return nil
}
