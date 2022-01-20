package processor

import (
	"context"
	"users-backup-aws-lambda/pkg/dto"
)

type (
	processorDummy struct{}
)

func NewProcessorDummy() *processorDummy {
	return &processorDummy{}
}

func (p *processorDummy) Process(ctx context.Context) (dto.Output, error) {
	return dto.Output{
		Mensaje: "Hola mundo",
	}, nil
}
