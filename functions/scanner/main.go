package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type S3Object struct {
	Bucket string
	Key    string
}

type Response struct {
	Status string
}

func scan(_ context.Context, toscan S3Object) (Response, error) {
	fmt.Printf("Received an event %+v\n", toscan)

	return Response{Status: "OK"}, nil

}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(scan)
}
