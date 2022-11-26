// Package notificationutil implements utility that should be used across all BriskSec products when it's required to
// send a notification to the central monitoring panel. This could include informational details, new detections, or
// any abnormal or error conditions that should be alerted to the administration team.
package notificationutil

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Building reusable connection to Amazon SQS.
var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = sqs.New(sess)

// Queue URL and MessageGroupId are coded into the library, so that other dependent applications do not have to worry
// about specifically configuration this information. If URL changes, a new version of the notificationutil should be
// released, and other projects should be updated to use the new version of the library.
var queueURL = "https://sqs.us-east-1.amazonaws.com/922655356721/aws_notification_bot_sqs.fifo"
var messageGroupId = "aws_notification_bot_sqs"

// Asia/Colombo timezone is used for timestamps within the messages.
var loc, _ = time.LoadLocation("Asia/Colombo")

// `Kind` denotes the type of notification getting sent. This will be used in determining the actions required at the
// receiving side.
type Kind uint64

const (
	Information Kind = iota
	UnexpectedError
	ExpectedError
	AbnormalCondition
)

// Returns the string representation of the notification Kind.
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

// Send a notification with the provided title, kind, message, and optionally an error. Error can be set to nil when
// the notification is not related to an error.
func Notify(title string, kind Kind, message string, err error) {
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
