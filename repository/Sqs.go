package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// NewSQS returns a new sns client for the passed in region
func newSQS(region, endpoint string) sqsiface.SQSAPI {
	cfg := aws.Config{
		Region: aws.String(region),
	}
	// if endpoint is not empty, we will use localstack
	if endpoint != "" {
		cfg.Endpoint = aws.String(endpoint)
	}

	sess := session.Must(session.NewSession(&cfg))

	cliSQS := sqs.New(sess)

	return cliSQS
}

func sendMessage(sqsClient sqsiface.SQSAPI, msg, queueURL string) (*sqs.SendMessageOutput, error) {
	sqsMessage := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(msg),
	}

	output, err := sqsClient.SendMessage(sqsMessage)
	if err != nil {
		return nil, fmt.Errorf("could not send message to queue %v: %v", queueURL, err)
	}

	return output, nil
}
