package notificationutil

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var svc = sqs.New(sess)

var queueURL = "https://sqs.us-east-1.amazonaws.com/922655356721/aws_notification_bot_sqs.fifo"
var messageGroupId = "aws_notification_bot_sqs"

var loc, _ = time.LoadLocation("Asia/Colombo")

type Kind uint64

const (
	Information Kind = iota
	UnexpectedError
	ExpectedError
	AbnormalCondition
)

func (i Kind) String() string {
	switch i {
	case Information:
		return "Information"
	case UnexpectedError:
		return "UnexpectedError"
	case ExpectedError:
		return "ExpectedError"
	case AbnormalCondition:
		return "AbnormalCondition"
	default:
		return "Unknown"
	}
}

func NotifyAbnormalCondition(title string, kind Kind, message string, err error) {
	details := ""
	if err != nil {
		details = fmt.Sprintf("%v", err)
		message = message + "\n\n" + details
	}

	_, aendErr := svc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(title),
			},
			"Kind": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("%d", kind)),
			},
			"Timestamp": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(time.Now().In(loc).String()),
			},
		},
		MessageBody:    aws.String(message),
		QueueUrl:       &queueURL,
		MessageGroupId: &messageGroupId,
	})

	if aendErr != nil {
		log.Printf("Error sending notification: %v \n", aendErr)
	}
}
