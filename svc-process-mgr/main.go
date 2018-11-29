package main

import (
	// "encoding/json"
	// "context"
	"flag"
	"lss_poc/svc-process-mgr/avro-kafka"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	goavro "gopkg.in/linkedin/goavro.v2"
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
	logger.Debug("KAFKA_BROKER_URI: ", topic)
	schemaRegistryURI = os.Getenv("SCHEMA_REGISTRY_URI")
	logger.Debug("SCHEMA_REGISTRY_URI: ", schemaRegistryURI)
}

func main() {
	flag.Parse()
	logger.Info("Servicing Process Manager started....")

	// go startIntakeListener()

	// go ExecuteProcess()
}

func startWorkPipeline() {

}

//
// func startIntakeListener() {
// 	// make a new reader that consumes
// 	r := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:   []string{brokerURI},
// 		Topic:     topic,
// 		Partition: 0,
// 		MinBytes:  10e3, // 10KB
// 		MaxBytes:  10e6, // 10MB
// 	})
// 	//Setting the offset ot -1 means to seek to the first offset. Setting the offset to -2 means to seek to the last offset.
// 	r.SetOffset(-2)
//
// 	for {
// 		m, err := r.ReadMessage(context.Background())
// 		if err != nil {
// 			logger.Error(err)
// 			break
// 		}
// 		logger.Printf("message at offset %d: %s = %s", m.Offset, string(m.Key), string(m.Value))
// 	}
//
// 	r.Close()
// }

func startSaramaConsumer() {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Consumer.Return.Errors = true

	producerUrl := strings.Split(brokerURI, ",")

	consumer, err := sarama.NewConsumer(producerUrl, config)
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	logger.Print("Connected to kafka broker")

	for m := range partitionConsumer.Messages() {
	}
}

func readAvroMessage(message *sarama.ConsumerMessage) {
	SchemaRegistryClient := avrokafka.NewSchemaRegistryClientWithRetries([]string{schemaRegistryURI}, 2)
	schemaCodec, err := SchemaRegistryClient.GetLatestSchema("loan")
	if err != nil {
		panic(err)
	}

	logger.Println("Schema: ", schemaCodec.Schema())
	codec, err := goavro.NewCodec(schemaCodec.Schema())
	if err != nil {
		panic(err)
	}

	encoded := []byte(message.Value)
	decoded, _, err := codec.NativeFromBinary(encoded)

	record := decoded.(*goavro.Record)
	log.Print("Record Name:", record.Name)
	log.Print("Record Fields:")
	for i, field := range record.Fields {
		log.Print("field", i, field.Name, ":", field.Datum)

		// if field.Name == "com.avro.kafka.golang.data" {
		//
		// 	data, _ := field.Datum.(string)
		// 	bytesText := []byte(data)
		//
		// 	var user User
		// 	json.Unmarshal(bytesText, &user)
		//
		// 	log.Print("raw data : ", data)
		// 	log.Print("id : ", user.ID)
		// 	log.Print("name : ", user.Name)
		//
		// }
	}
}
