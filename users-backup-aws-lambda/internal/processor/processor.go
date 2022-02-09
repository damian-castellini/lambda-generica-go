package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"users-backup-aws-lambda/pkg/dto"
)

type (
	secretInterface interface {
		GetDBSecret() string
	}
	svcInterface interface {
		GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	}

	processor struct {
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

type Secret struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Db       string `json:"db"`
	Server   string `json:"server"`
}

func NewProcessor(s secretInterface, s3 svcInterface) *processor {
	return &processor{dbSecret: s, svc: s3}
}

func (p *processor) Process(ctx context.Context, s3Event events.S3Event) (dto.Output, error) {
	var finalMetrics []Metrics
	var secret Secret
	resp := p.dbSecret.GetDBSecret()
	errorParsingSecret := json.Unmarshal([]byte(resp), &secret)

	if errorParsingSecret != nil {
		fmt.Println(errorParsingSecret)
	} else {
		fmt.Println(secret.Db, secret.Ip)
	}

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
			fmt.Printf("\nprocessing json %s", jsonBody)
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
