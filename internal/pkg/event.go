package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const DelaySeconds = 1
const MaxNumberOfMessages = 10
const VisibilityTimeoutSeconds = 5

type Session struct {
	sess     *session.Session
	svc      *sqs.SQS
	QueueURL *string
}

func NewSession() (*Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return nil, err
	}

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	svc := sqs.New(sess)
	queue, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("gophermart"),
	})
	if err != nil {
		return nil, err
	}
	return &Session{
		sess:     sess,
		svc:      svc,
		QueueURL: queue.QueueUrl,
	}, err
}

func (s *Session) SendMessage(message string) (*sqs.SendMessageOutput, error) {
	return s.svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(DelaySeconds),
		MessageBody:  aws.String(message),
		QueueUrl:     s.QueueURL,
	})
}

func (s *Session) ReceiveMessages() (result []*sqs.Message, err error) {
	output, err := s.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            s.QueueURL,
		MaxNumberOfMessages: aws.Int64(MaxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(VisibilityTimeoutSeconds),
	})
	if err != nil {
		return result, nil
	}
	return append(result, output.Messages...), nil
}

func (s *Session) DeleteMessage(message *sqs.Message) error {
	_, err := s.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      s.QueueURL,
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}
