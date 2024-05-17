package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oklog/ulid/v2"
	"github.com/sethvargo/go-envconfig"
	"log"
	"time"
)

type PresignResponse struct {
	URL    string
	Method string
}

type PreSigner struct {
	Bucket string `env:"BUCKET, required"`
}

func (p PreSigner) signedUrl(ctx context.Context) (*PresignResponse, error) {

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(sdkConfig)
	presignClient := s3.NewPresignClient(s3Client)
	presign, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(ulid.Make().String()),
	}, s3.WithPresignExpires(time.Minute*10))

	if err != nil {
		return nil, err
	}

	return &PresignResponse{
		URL:    presign.URL,
		Method: presign.Method,
	}, nil
}

func main() {
	ctx := context.Background()
	var app PreSigner
	if err := envconfig.Process(ctx, &app); err != nil {
		log.Fatal(err)
	}
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(app.signedUrl)
}
