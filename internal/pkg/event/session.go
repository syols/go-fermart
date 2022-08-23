package event

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Session struct {
	sess     *session.Session
	svc      *sqs.SQS
	QueueURL *string
}

func NewSession() (*Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("AKIA52TFFBZFP3TOQRG2", "QB2JzDizt62e1NQR+Qvc77D4fUH3M1cx9ofmnIK3", ""),
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

func (s Session) SendMessage(message string) (*sqs.SendMessageOutput, error) {
	return s.svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message),
		QueueUrl:     s.QueueURL,
	})
}

func (s Session) ReceiveMessages() (result []*sqs.Message, err error) {
	output, err := s.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            s.QueueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(5),
	})
	if err != nil {
		return result, nil
	}
	return append(result, output.Messages...), nil
}

func (s Session) DeleteMessage(message *sqs.Message) error {
	log.Print("DeleteMessage")
	_, err := s.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      s.QueueURL,
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}
