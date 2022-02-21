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

	databaseInterface interface {
		Open(connectionString string)
		MigrateUser(userToInsert string) (int64, error)
	}

	processor struct {
		dbSecret     secretInterface
		svc          svcInterface
		dbConnection databaseInterface
	}
)

func NewProcessor(s secretInterface, s3 svcInterface, db databaseInterface) *processor {
	return &processor{dbSecret: s, svc: s3, dbConnection: db}
}

func (p *processor) Process(ctx context.Context, s3Event events.S3Event) (dto.Output, error) {
	var finalMetrics []dto.Metrics
	resp := p.dbSecret.GetDBSecret()
	p.dbConnection.Open(resp)

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
			var metrics dto.Metrics
			err = json.Unmarshal([]byte(jsonBody), &metrics)
			if err != nil {
				fmt.Print(err)
			} else {
				fmt.Printf(metrics.Message)
				insertedId, err := p.dbConnection.MigrateUser(metrics.Message)
				if insertedId == -1 {
					fmt.Println(dto.ERROR_INSERT, err)
				} else {
					fmt.Println(dto.INSERT_SUCCESSFUL)
				}
				finalMetrics = append(finalMetrics, metrics)
			}
		}

	}
	return dto.Output{
		Message: dto.STARTING_PROCESS,
	}, nil
}
