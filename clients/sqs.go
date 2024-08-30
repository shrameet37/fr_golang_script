package clients

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"

	"face_management/logger"
)

var (
	ToSqsChEventAggregatorQueue = make(chan string, 200)
)

func SqsMessageSender(queueURL string, sqsChan <-chan string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		logger.Log.Fatal("Failed to create AWS session:", zap.Error(err))
	}

	// Create the SQS client using the AWS session
	svc := sqs.New(sess)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for message := range sqsChan {
			sendSQSMessage(svc, queueURL, message)
		}
	}()

	logger.Log.Info("SqsMessageSender, Started SQS Message Sender Routine")
	wg.Wait()
}

func sendSQSMessage(svc *sqs.SQS, queueURL string, message string) {

	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Recovered from panic", zap.Any("recover", r))
		}
	}()

	msgHash := getMD5Hash(message)

	logger.Log.Info("SQS MESSAGE SENDER", zap.String("msgHash", msgHash), zap.String("message", message))

	sendMessageOutput, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": {
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
		},
		MessageBody: aws.String(message),
		QueueUrl:    &queueURL,
	})

	if err != nil {
		logger.Log.Error("SQS MESSAGE", zap.Error(err))
		return
	}

	logger.Log.Info("SQS MESSAGE SENT", zap.String("res", sendMessageOutput.String()))
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
