package handler

import (
	"context"

	"users-backup-aws-lambda/pkg/dto"
)

type (
	processor interface {
		Process(context.Context) (dto.Output, error)
	}

	LambdaHandler struct {
		processor processor
	}
)

func NewLambdaHandler(p processor) *LambdaHandler {
	return &LambdaHandler{processor: p}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context) (dto.Output, error) {

	return h.processor.Process(ctx)
}
