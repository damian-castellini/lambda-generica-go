package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"

	"users-backup-aws-lambda/pkg/dto"
)

type (
	processor interface {
		Process(context.Context, events.S3Event) (dto.Output, error)
	}

	LambdaHandler struct {
		processor processor
	}
)

func NewLambdaHandler(p processor) *LambdaHandler {
	return &LambdaHandler{processor: p}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, s3Event events.S3Event) (dto.Output, error) {
	return h.processor.Process(ctx, s3Event)
}
