package main

import (
   "fmt"
   "flag"
   "github.com/aws/aws-sdk-go/aws"
   aws_session "github.com/aws/aws-sdk-go/aws/session"
   "github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
    var timeout int64
    flag.Int64Var(&timeout, "t", 20, "(Optional) Timeout in seconds for long polling")
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
   var receiptHandles []string 
   for _, message := range receivedMessages {
       receiptHandles = append(receiptHandles, *message.ReceiptHandle)
   }
   fmt.Println(receiptHandles) 
}
