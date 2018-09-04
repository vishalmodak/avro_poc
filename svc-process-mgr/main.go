package main

import (
	// "encoding/json"
	"context"
	"flag"
	"os"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

var (
	topic     string
	brokerURI string
	logLevel  string
)

func init() {
	gotenv.Load()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	topic = os.Getenv("TOPIC_NAME")
	log.Debug("TOPIC_NAME: ", topic)
	brokerURI = os.Getenv("KAFKA_BROKER_URI")
	log.Debug("KAFKA_BROKER_URI: ", topic)
}

func main() {
	flag.Parse()
	log.Info("Servicing Process Manager started....")

	startIntakeListener()
}

func startIntakeListener() {
	// make a new reader that consumes
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerURI},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	//Setting the offset ot -1 means to seek to the first offset. Setting the offset to -2 means to seek to the last offset.
	r.SetOffset(-2)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Error(err)
			break
		}
		log.Printf("message at offset %d: %s = %s", m.Offset, string(m.Key), string(m.Value))
	}

	r.Close()
}