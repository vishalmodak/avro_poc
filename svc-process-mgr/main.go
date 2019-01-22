package main

import (
	// "encoding/json"
	"bytes"
	"flag"
	"github.com/Shopify/sarama"
	// avro "github.com/elodina/go-avro"
	kavro "github.com/elodina/go-kafka-avro"
	goavro "github.com/linkedin/goavro"
	// "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	gavro "gopkg.in/avro.v0"
	"io/ioutil"
	// lg "log"
	// "lss_poc/svc-process-mgr/avro-kafka"
	"os"
	"strings"
)

var (
	loanTopic         string
	paymentTopic      string
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

	loanTopic = os.Getenv("LOAN_TOPIC_NAME")
	logger.Debug("LOAN_TOPIC_NAME: ", loanTopic)
	paymentTopic = os.Getenv("PAYMENT_TOPIC_NAME")
	logger.Debug("PAYMENT_TOPIC_NAME: ", paymentTopic)
	brokerURI = os.Getenv("KAFKA_BROKER_URI")
	logger.Debug("KAFKA_BROKER_URI: ", brokerURI)
	schemaRegistryURI = os.Getenv("SCHEMA_REGISTRY_URI")
	logger.Debug("SCHEMA_REGISTRY_URI: ", schemaRegistryURI)
}

func main() {
	flag.Parse()
	logger.Info("Servicing Process Manager started....")

	go startLoanConsumer()

	// go ExecuteProcess()

	startPaymentConsumer()
}

func startLoanConsumer() {
	config := sarama.NewConfig()
	// config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Consumer.Return.Errors = true

	//verbose debugging (comment this line to disabled verbose sarama logging)
	// sarama.Logger = lg.New(os.Stdout, "[sarama] ", lg.LstdFlags)

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

	partitions, err := consumer.Partitions(loanTopic)
	if err != nil {
		logger.Fatalf("Failed to get partitions: %s", err)
	}
	logger.Printf("%d partitions for topic: %s", len(partitions), loanTopic)

	partitionConsumer, err := consumer.ConsumePartition(loanTopic, 0, sarama.OffsetNewest)
	if err != nil {
		logger.Errorf("Failed to consume partition: %s", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	logger.Print("Started Loan Consumer...")

	for m := range partitionConsumer.Messages() {
		readAvroLoan(m.Value)
	}
}

func startPaymentConsumer() {
	config := sarama.NewConfig()
	// config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Consumer.Return.Errors = true

	//verbose debugging (comment this line to disabled verbose sarama logging)
	// sarama.Logger = lg.New(os.Stdout, "[sarama] ", lg.LstdFlags)

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

	partitions, err := consumer.Partitions(paymentTopic)
	if err != nil {
		logger.Fatalf("Failed to get partitions: %s", err)
	}
	logger.Printf("%d partitions for topic: %s", len(partitions), paymentTopic)

	partitionConsumer, err := consumer.ConsumePartition(paymentTopic, 0, sarama.OffsetNewest)
	if err != nil {
		logger.Errorf("Failed to consume partition: %s", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	logger.Print("Started Payment Consumer....")

	for m := range partitionConsumer.Messages() {
		readAvroMessage4(m.Value)
	}
}

//uses linkedin/avro
//doesnt work, fails decoding from binary to native
//Actual Output
//   map[string]interface {}{"loanNumber":"", "sourceAccountNumber":"", "currentOwner":""}
func readAvroMessage(encodedMsg []byte) {
	// SchemaRegistryClient := avrokafka.NewSchemaRegistryClientWithRetries([]string{schemaRegistryURI}, 2)
	// schemaCodec, err := SchemaRegistryClient.GetLatestSchema("loan")
	// if err != nil {
	// 	logger.Errorf("Failure fetching latest schema for loan: %s", err)
	// }
	//
	// logger.Debug("Schema: ", schemaCodec.Schema())

	schema, err := ioutil.ReadFile("loan.avsc")
	logger.Println(string(schema))
	codec, err := goavro.NewCodec(string(schema))
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

//uses github.com/elodina/go-kafka-avro
//doesnt work
//Actual output
//    main.Loan{LoanNumber:"current", SourceAccountNumber:"", CurrentOwner:""}
//    &main.Payment{LoanNumber:"", AmountInCents:0, SourceAccountNumber:"", SourcePaymentNumber:"", SourceObligationNumber:""}
func readAvroMessage2(encodedMsg []byte) {
	decoder := kavro.NewKafkaAvroDecoder(schemaRegistryURI)
	// decoded, err := decoder.Decode(encodedMsg)
	// if err != nil {
	// 	logger.Errorf("Failure decoding avro: %s", err)
	// }
	// decodedRecord, ok := decoded.(*avro.GenericRecord)
	// if !ok {
	// 	logger.Errorf("Failure casting to GenericRecord")
	// }
	// logger.Printf("%v", decodedRecord)

	payment := new(Payment)
	err := decoder.DecodeSpecific(encodedMsg, payment)
	if err != nil {
		logger.Errorf("Failure decoding avro: %s", err)
	}
	logger.Printf("%#v", payment)
}

//uses https://github.com/actgardner/gogen-avro
//doesnt work
//Actual output
//   &main.Loan{LoanNumber:"", SourceAccountNumber:""}
func readAvroMessage3(encodedMsg []byte) {
	r := bytes.NewReader(encodedMsg)
	logger.Info(encodedMsg)
	loan, err := DeserializeLoan(r)
	if err != nil {
		logger.Errorf("Failure during loan deserialization: %s", err)
	}
	logger.Printf("%#v", loan)

}

//uses https://github.com/go-avro/avro for decoding
//doesnt work
//Actual output
//   &main.Loan{LoanNumber:"current", SourceAccountNumber:""}
//   &main.Payment{LoanNumber:"", AmountInCents:17, SourceAccountNumber:"aid\":true,\"datePaid\":\"2015-12-23\",
//     \"loanNumber\":\"2015CA16", SourcePaymentNumber:"", SourceObligationNumber:""}
//   &main.Payment{Paid:false, DatePaid:"paid\":true,\"dateP", LoanNumber:"", AmountInCents:-53, SourceAccountNumber:"\":\"2015-12-23\",
//     \"loanNumber\":\"2015CA169772974\",\"amo", SourcePaymentNumber:"", SourceObligationNumber:"tInCents\":12151,
//     \"sourceAccountNumber\":\"8601860\",\"source"}
func readAvroMessage4(encodedMsg []byte) {

	// payment_list_schema, err := gavro.ParseSchemaFile("payment_list.avsc")
	// if err != nil {
	// 	// Should not happen if the schema is valid
	// 	logger.Errorf("Failure parsing schema: %s", err)
	// }
	payment_schema, err := gavro.ParseSchemaFile("loan.avsc")
	if err != nil {
		// Should not happen if the schema is valid
		logger.Errorf("Failure parsing schema: %s", err)
	}
	reader := gavro.NewSpecificDatumReader()
	// SetSchema must be called before calling Read
	reader.SetSchema(payment_schema)
	// Create a new Decoder with a given buffer
	decoder := gavro.NewBinaryDecoder(encodedMsg)

	decodedRecord := new(Payment)

	// Read data into a given record with a given Decoder.
	err = reader.Read(decodedRecord, decoder)
	if err != nil {
		logger.Errorf("Failure decoding avro: %s", err)
	}

	logger.Printf("%#v\n", decodedRecord)
}

func readAvroLoan(encodedMsg []byte) {

	// payment_list_schema, err := gavro.ParseSchemaFile("payment_list.avsc")
	// if err != nil {
	// 	// Should not happen if the schema is valid
	// 	logger.Errorf("Failure parsing schema: %s", err)
	// }
	loan_schema, err := gavro.ParseSchemaFile("loan.avsc")
	if err != nil {
		// Should not happen if the schema is valid
		logger.Errorf("Failure parsing schema: %s", err)
	}
	reader := gavro.NewSpecificDatumReader()
	// SetSchema must be called before calling Read
	reader.SetSchema(loan_schema)
	// Create a new Decoder with a given buffer
	decoder := gavro.NewBinaryDecoder(encodedMsg)

	decodedRecord := new(Loan)

	// Read data into a given record with a given Decoder.
	err = reader.Read(decodedRecord, decoder)
	if err != nil {
		logger.Errorf("Failure decoding avro: %s", err)
	}

	logger.Printf("%#v\n", decodedRecord)
}
