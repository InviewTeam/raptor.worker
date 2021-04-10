package kafka

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
)

func NewProducer(brokers []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionNone
	config.Producer.MaxMessageBytes = 100000000

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	//Relay incoming signals to channel 'c'
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {
		sig := <-c
		log.Println("Got signal:", sig)

		if err := producer.Close(); err != nil {
			log.Fatal("Error closing async producer:", err)
		}

		log.Println("Async Producer closed.")
		os.Exit(1)
	}()

	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write message to topic:", err)
		}
	}()

	return producer, nil
}
