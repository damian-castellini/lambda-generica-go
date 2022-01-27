package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"log"
	"strings"
	"users-backup-aws-lambda/pkg/dto"
)

type (
	secretInterface interface {
		GetDBSecret() string
	}
	svcInterface interface {
		GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	}

	processorDummy struct {
		dbSecret secretInterface
		svc      svcInterface
	}
)

type Metrics struct {
	Type           string
	MessageId      string
	TopicArn       string
	Subject        string
	Message        string
	Timestamp      string
	UnsubscribeURL string
}

func NewProcessorDummy(s secretInterface, s3 svcInterface) *processorDummy {
	return &processorDummy{dbSecret: s, svc: s3}
}

func (p *processorDummy) Process(ctx context.Context, s3Event events.S3Event) (dto.Output, error) {
	var finalMetrics []Metrics
	var resp = p.dbSecret.GetDBSecret()
	fmt.Println(resp)
	for _, record := range s3Event.Records {
		var bucketName = record.S3.Bucket.Name
		var fileKey = record.S3.Object.Key
		fmt.Printf("Bucket = %s, Key = %s", bucketName, fileKey)
		requestInput := &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileKey),
		}
		result, err := p.svc.GetObject(requestInput)
		if err != nil {
			log.Print(err)
		}

		defer result.Body.Close()
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			log.Print(err)
		}

		bodyString := fmt.Sprintf("%s", body)
		var splitBody = strings.Split(bodyString, "\n")
		for _, jsonBody := range splitBody {
			fmt.Printf(jsonBody)
			fmt.Printf("processing json %s", jsonBody)
			var metrics Metrics
			err = json.Unmarshal([]byte(jsonBody), &metrics)
			if err != nil {
				fmt.Print(err)
			} else {
				fmt.Printf(metrics.Message)
				finalMetrics = append(finalMetrics, metrics)
			}
		}

	}
	return dto.Output{
		Mensaje: "Hola mundo",
	}, nil
}
