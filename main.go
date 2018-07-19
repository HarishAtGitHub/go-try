package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strconv"
)

func main() {
	var timeout int64 = 20
	fmt.Println(timeout)
	sess := aws_session.Must(aws_session.NewSession())
	svc := sqs.New(sess)
	// get gueue url from queue name and account ID
	qName :=  "mammoth-config-sqs"
	qOwnerAWSAccountId := "603514901691"
	queueUrlOutput, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &qName,
		QueueOwnerAWSAccountId: &qOwnerAWSAccountId,
	})
	qURL := *queueUrlOutput.QueueUrl

	// get queue messages in bulk
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(4),
		WaitTimeSeconds:     aws.Int64(timeout),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}
	//fmt.Println(result.Messages)

	// delete messages in bulk
	receivedMessages := result.Messages
	var entries []*sqs.DeleteMessageBatchRequestEntry
	for _, message := range receivedMessages {
		// send message to hec
		// if success full
		entry := &sqs.DeleteMessageBatchRequestEntry{Id: message.MessageId, ReceiptHandle: message.ReceiptHandle}
		entries = append(entries, entry)
	}

	deleteParams := &sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(qURL),
	}

	deleteResp, err := svc.DeleteMessageBatch(deleteParams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(deleteResp)

	// get size of queue
	attributes, err := svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: &qURL,
		AttributeNames: []*string{
			aws.String("All"),
		},
	})
	for attrib, _ := range attributes.Attributes {
		prop := attributes.Attributes[attrib]
		i, _ := strconv.Atoi(*prop)
		fmt.Println(attrib, i)
	}
	fmt.Println(*attributes.Attributes["ApproximateNumberOfMessages"])
}
