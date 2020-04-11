package backends

import (
	"async-queue/broker"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"sync"
)

var (
	singleton  sync.Once
	sqsService *sqs.SQS
)

const (
	pollMessageCount = 1
	queueUrl         = "https://sqs.us-west-1.amazonaws.com/<queue-id>/<queue-name>"
	region           = "us-west-1"
	credPath         = "/Users/<user-name>/.aws/credentials"
	credProfile      = "default"
	maxRetries       = 5
)

type SQSBackend struct {
	region      string
	credPath    string
	credProfile string
	queueUrl    string
}

func NewSQSService() broker.QueueBackendManager {
	return &SQSBackend{
		region:      region,
		credPath:    credPath,
		credProfile: credProfile,
		queueUrl:    queueUrl,
	}
}

func (s *SQSBackend) GetConnection() interface{} {
	singleton.Do(func() {
		sess, _ := session.NewSession(&aws.Config{
			Region:      aws.String(s.region),
			Credentials: credentials.NewSharedCredentials(s.credPath, s.credProfile),
			MaxRetries:  aws.Int(maxRetries),
		})
		sqsService = sqs.New(sess)
	})
	return sqsService
}

func (s *SQSBackend) EnqueueTask(message string) error {
	svc := s.GetConnection().(*sqs.SQS)
	sent, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(string(message)),
		QueueUrl:     &s.queueUrl,
		DelaySeconds: aws.Int64(3),
	})
	if err != nil {
		return err
	}
	log.Printf("Sent message %v", sent)
	return nil
}

func (s *SQSBackend) DequeueTask() map[*string]*string {
	svc := s.GetConnection().(*sqs.SQS)
	response, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueUrl),
		MaxNumberOfMessages: aws.Int64(pollMessageCount),
		VisibilityTimeout:   aws.Int64(30),
		//WaitTimeSeconds:     aws.Int64(10), //Long polling for 10 seconds
	})
	if err != nil {
		log.Printf("Error receiving message from sqs %s", err)
	}
	log.Printf("Received response from sqs -- %+v", response)
	msgLen := len(response.Messages)
	messageIdentifiers := map[*string]*string{}
	if msgLen > 0 {
		for _, msg := range response.Messages {
			messageIdentifiers[msg.ReceiptHandle] = msg.Body
		}
	}
	return messageIdentifiers
}

func (s *SQSBackend) DeleteTask(receiptHandle *string) {
	_, err := sqsService.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueUrl),
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		log.Printf("Error deleting message %s", err)
	}
	log.Printf("Message ID: %+v has beed deleted", receiptHandle)
}
