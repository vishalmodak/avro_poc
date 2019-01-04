package main

import (
	// "encoding/json"
	"context"
	"flag"
	"github.com/Shopify/sarama"
	avro "github.com/elodina/go-avro"
	kavro "github.com/elodina/go-kafka-avro"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	goavro "gopkg.in/linkedin/goavro.v2"
	lg "log"
	"lss_poc/svc-process-mgr/avro-kafka"
	"os"
	"strings"
)

var (
	topic             string
	brokerURI         string
	logLevel          string
	logger            *log.Logger
	schemaRegistryURI string
)

func init() {
	gotenv.Load()
	logger = log.New()
	// log.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	topic = os.Getenv("TOPIC_NAME")
	logger.Debug("TOPIC_NAME: ", topic)
	brokerURI = os.Getenv("KAFKA_BROKER_URI")
	logger.Debug("KAFKA_BROKER_URI: ", brokerURI)
	schemaRegistryURI = os.Getenv("SCHEMA_REGISTRY_URI")
	logger.Debug("SCHEMA_REGISTRY_URI: ", schemaRegistryURI)
}

func main() {
	flag.Parse()
	logger.Info("Servicing Process Manager started....")

	// startIntakeListener()

	// go ExecuteProcess()

	startSaramaConsumer()
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
			logger.Error(err)
			break
		}
		logger.Printf("message at offset %d: %s = %s", m.Offset, string(m.Key), string(m.Value))
		readAvroMessage2(m.Value)
	}

	r.Close()
}

func startSaramaConsumer() {
	config := sarama.NewConfig()
	// config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Consumer.Return.Errors = true

	//verbose debugging (comment this line to disabled verbose sarama logging)
	sarama.Logger = lg.New(os.Stdout, "[sarama] ", lg.LstdFlags)

	var producerURL []string
	if strings.Contains(brokerURI, ",") {
		producerURL = strings.Split(brokerURI, ",")
	} else {
		producerURL = append(producerURL, brokerURI)
	}

	consumer, err := sarama.NewConsumer(producerURL, config)
	if err != nil {
		logger.Fatalf("Failed to start consumer: %s", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		logger.Fatalf("Failed to get partitions: %s", err)
	}
	logger.Printf("%d partitions for topic: %s", len(partitions), topic)

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logger.Errorf("Failed to consume partition: %s", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	logger.Print("Connected to kafka broker")

	for m := range partitionConsumer.Messages() {
		readAvroMessage2(m.Value)
	}
}

func readAvroMessage(encodedMsg []byte) {
	SchemaRegistryClient := avrokafka.NewSchemaRegistryClientWithRetries([]string{schemaRegistryURI}, 2)
	schemaCodec, err := SchemaRegistryClient.GetLatestSchema("loan")
	if err != nil {
		logger.Errorf("Failure fetching latest schema for loan: %s", err)
	}

	logger.Debug("Schema: ", schemaCodec.Schema())

	codec, err := goavro.NewCodec(schemaCodec.Schema())
	if err != nil {
		logger.Errorf("Failure reading code from schema: %s", err)
	}

	decoded, original, err := codec.NativeFromBinary(encodedMsg)
	if err != nil {
		logger.Errorf("Failure decoding binary data: %s", err)
		logger.Errorf("Original: %s", original)
	}

	logger.Printf("%#v", decoded)
}

func readAvroMessage2(encodedMsg []byte) {
	decoder := kavro.NewKafkaAvroDecoder(schemaRegistryURI)
	decoded, err := decoder.Decode(encodedMsg)
	if err != nil {
		logger.Errorf("Failure decoding avro: %s", err)
	}
	decodedRecord, ok := decoded.(*avro.GenericRecord)
	if !ok {
		logger.Errorf("Failure casting to GenericRecord")
	}
	logger.Printf("%v", decodedRecord)

	// var loan *Loan
	// err = decoder.DecodeSpecific(encodedMsg, &loan)
	// if err != nil {
	// 	logger.Errorf("Failure decoding avro: %s", err)
	// }
	// logger.Printf("%#v", *loan)
}
