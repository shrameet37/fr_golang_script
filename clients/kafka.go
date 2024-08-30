package clients

import (
	"face_management/logger"
	"face_management/misc"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type callbackFunctionWithMsg func(*kafka.Message) error

type ToKafkaMessage struct {
	Topic string
	Key   string
	Value []byte
}

var (
	ToKafkaChResourceAckEvent          = make(chan ToKafkaMessage)
	ToKafkaChToIotEngineCommandEvent   = make(chan ToKafkaMessage)
	ToKafkaChToAcaasPermissionAckEvent = make(chan ToKafkaMessage)
)

func KafkaConsumer(kafkaConsumerGroup string, kafkaTopicName string, callbackFunction callbackFunctionWithMsg) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":           kafkaConsumerGroup,
		"auto.offset.reset":  "latest",
		"enable.auto.commit": false,
		"sasl.mechanisms":    "SCRAM-SHA-512",
		"security.protocol":  "sasl_ssl",
		"sasl.username":      os.Getenv("KAFKA_USERNAME"),
		"sasl.password":      os.Getenv("KAFKA_PASSWORD"),
	})

	if err != nil {
		errorMsg := fmt.Sprintf("Kafka consumer connection error.Topic:%s,Error:%s!\n", kafkaTopicName, err.Error())
		logger.Log.Error(errorMsg)
		misc.ProcessError(9, errorMsg, nil)
		panic(err)
	}

	err = c.SubscribeTopics([]string{kafkaTopicName}, nil)

	if err != nil {
		errorMsg := fmt.Sprintf("Kafka consumer subscribe error.Topic:%s,Error:%s!\n", kafkaTopicName, err.Error())
		logger.Log.Error(errorMsg)
		misc.ProcessError(9, errorMsg, nil)
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			err = callbackFunction(msg)
			if err == nil {
				c.CommitMessage(msg)
			} else {
				errorMsg := fmt.Sprintf("Kafka consumer commit error.Error:%s,Message:%s!\n", err.Error(), string(msg.Value))
				logger.Log.Error(errorMsg)
				misc.ProcessError(9, errorMsg, msg.Value)
				c.Close()
				return
			}
		} else {
			//Here I am making the service panic and restart whenever consumer read error occurs
			errorMsg := fmt.Sprintf("Kafka consumer read error.Error:%s!\n", err.Error())
			logger.Log.Error(errorMsg)
			misc.ProcessError(9, errorMsg, nil)
			panic(err)
			//c.Close()
			//time.Sleep(5 * time.Second)
			//go KafkaConsumer(kafkaConsumerGroup, kafkaTopicName, callbackFunction)
			//return
		}
	}

}

func KafkaProducer(producerChannel <-chan ToKafkaMessage) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"sasl.mechanisms":   "SCRAM-SHA-512",
		"security.protocol": "sasl_ssl",
		"sasl.username":     os.Getenv("KAFKA_USERNAME"),
		"sasl.password":     os.Getenv("KAFKA_PASSWORD"),
	})

	if err != nil {
		errorMsg := fmt.Sprintf("Kafka producer connection error.Error:%s!\n", err.Error())
		logger.Log.Error(errorMsg)
		misc.ProcessError(9, errorMsg, nil)
		panic(err)
	}

	deliveryChan := make(chan kafka.Event)

	for {

		message := <-producerChannel

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &(message.Topic), Partition: kafka.PartitionAny},
			Key:            []byte(message.Key),
			Value:          []byte(message.Value),
		}, deliveryChan)

		kafkaEvent := <-deliveryChan
		kafkaMessage := kafkaEvent.(*kafka.Message)

		if kafkaMessage.TopicPartition.Partition == -1 {
			errorMsg := fmt.Sprintf("Kafka producer or topic partition error.Topic:%s,Message:%s!\n", *kafkaMessage.TopicPartition.Topic, string(kafkaMessage.Value))
			logger.Log.Error(errorMsg)
			misc.ProcessError(9, errorMsg, kafkaMessage.Value)
		}

		if err != nil {
			errorMsg := fmt.Sprintf("Kafka producer produce error.Error:%s!\n", err.Error())
			logger.Log.Error(errorMsg)
			misc.ProcessError(9, errorMsg, nil)
		}

	}

}
