package main

import (
	// "encoding/json"
	"bytes"
	"flag"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
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
		readLoanMessage(m.Value[5:])
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
		readPayments(m.Value[5:])
	}
}

//uses https://github.com/actgardner/gogen-avro
//doesnt work
//Actual output
//   &main.Loan{LoanNumber:"", SourceAccountNumber:""}
func readPayments(encodedMsg []byte) {
	r := bytes.NewReader(encodedMsg)
	paymentList, err := DeserializePayment_list(r)
	if err != nil {
		logger.Errorf("Failure during paymentList deserialization: %s", err)
	}
	logger.Printf("%+v", paymentList)

}

//works
//Output
// &main.Payment{Paid:true, DatePaid:"2015-12-23", LoanNumber:"2015CA169772974",
//               AmountInCents:12151, SourceAccountNumber:"8601860", SourcePaymentNumber:"2012323",
//               SourceObligationNumber:"11544267"}
func readPay(encodedMsg []byte) {
	r := bytes.NewReader(encodedMsg)
	logger.Println(string(encodedMsg))
	payment, err := DeserializePayment(r)
	if err != nil {
		logger.Errorf("Failure during payment deserialization: %s", err)
	}
	logger.Printf("%+v", payment)

}

//Works
func readLoanMessage(encodedMsg []byte) {
	r := bytes.NewReader(encodedMsg)
	loan, err := DeserializeLoan(r)
	if err != nil {
		logger.Errorf("Failure during loan deserialization: %s", err)
	}
	logger.Printf("%#v", loan)
}
