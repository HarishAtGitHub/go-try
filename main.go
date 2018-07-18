package main

import (
   "fmt"
   "github.com/aws/aws-sdk-go/aws"
   aws_session "github.com/aws/aws-sdk-go/aws/session"
   "github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
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

    // get queue message
    result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
        AttributeNames: []*string{
            aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
        },
        MessageAttributeNames: []*string{
            aws.String(sqs.QueueAttributeNameAll),
        },
        QueueUrl:            &qURL,
        MaxNumberOfMessages: aws.Int64(10),
        VisibilityTimeout:   aws.Int64(36000),  // 10 hours
        WaitTimeSeconds:     aws.Int64(0),
    })

    if err != nil {
        fmt.Println("Error", err)
        return
    }

    if len(result.Messages) == 0 {
        fmt.Println("Received no messages")
        return
    }
    fmt.Println(result)

   // get queue bulk messages
}
