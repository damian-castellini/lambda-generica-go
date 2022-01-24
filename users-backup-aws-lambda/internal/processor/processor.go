package processor

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"users-backup-aws-lambda/pkg/dto"
)

type (
	secretInterface interface {
		GetDBSecret() string
	}
	processorDummy struct {
		dbSecret secretInterface
	}
)

func NewProcessorDummy(s secretInterface) *processorDummy {
	return &processorDummy{dbSecret: s}
}

func (p *processorDummy) Process(ctx context.Context, s3Event events.S3Event) (dto.Output, error) {
	var resp = p.dbSecret.GetDBSecret()
	fmt.Println(resp)
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
	}
	return dto.Output{
		Mensaje: "Hola mundo",
	}, nil
}
