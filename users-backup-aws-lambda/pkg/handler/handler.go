package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"

	"users-backup-aws-lambda/pkg/dto"
)

type (
	processor interface {
		Process(context.Context, events.S3Event) (dto.Output, error)
	}
	secretInterface interface {
		GetDBSecret() string
	}

	LambdaHandler struct {
		secret    secretInterface
		processor processor
	}
)

func NewLambdaHandler(p processor, s secretInterface) *LambdaHandler {
	return &LambdaHandler{processor: p, secret: s}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, s3Event events.S3Event) (dto.Output, error) {
	var secretString = h.secret.GetDBSecret()
	fmt.Println(secretString)
	return h.processor.Process(ctx, s3Event)
}
