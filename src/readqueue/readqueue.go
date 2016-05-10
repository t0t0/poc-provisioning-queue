package main

import (
    "fmt"
    "os"
    "log"


    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	getMessage()
}

func getMessage() {
	svc := sqs.New(session.New())

params := &sqs.ReceiveMessageInput{
	QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/659527370395/provision-queue"), // Required
	AttributeNames: []*string{
		aws.String("QueueAttributeName"), // Required
		// More values...
	},
	MaxNumberOfMessages: aws.Int64(10),
	MessageAttributeNames: []*string{
		aws.String(""), // Required
		// More values...
	},
}


resp, err := svc.ReceiveMessage(params)

if err != nil {
	// Print the error, cast err to awserr.Error to get the Code and
	// Message from an error.
	fmt.Println(err.Error())
	return
}

l := resp.Messages
for i := 0; i < len(l); i ++ {
    v := l[i].Body
    handler := l[i].ReceiptHandle
    fmt.Println(*v)
    getFile(*v , *handler)
}

}

func getFile(key string, handler string) {
path := "tmp/"
location := path + key
file, err := os.Create(location)
if err != nil {
 log.Fatal("Failed to create file", err)
}
defer file.Close()
downloader := s3manager.NewDownloader(session.New(&aws.Config{Region:
aws.String("us-east-1")}))
numBytes, err := downloader.Download(file,
 &s3.GetObjectInput{
 Bucket: aws.String("poc-provision"),
 Key: aws.String(key),
 })
if err != nil {
 fmt.Println("Failed to download file", err)
 return
}
fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")
deleteMessage(handler, "https://sqs.us-east-1.amazonaws.com/659527370395/provision-queue")
	
}

func deleteMessage(r string, q string) {

	svc := sqs.New(session.New())

	params := &sqs.DeleteMessageInput{
	QueueUrl:      aws.String(q), // Required
	ReceiptHandle: aws.String(r), // Required
}
resp, err := svc.DeleteMessage(params)

if err != nil {
	// Print the error, cast err to awserr.Error to get the Code and
	// Message from an error.
	fmt.Println(err.Error())
	return
}

// Pretty-print the response data.
fmt.Println(resp)
	
}