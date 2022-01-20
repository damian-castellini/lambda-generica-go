package internal

import (
	"github.com/Bancar/uala-go-platform-product-dependencies/pkg/errors"
	_ "github.com/aws/aws-sdk-go/service/dynamodb"
	_ "users-backup-aws-lambda/internal/aws"
	"users-backup-aws-lambda/internal/processor"
	"users-backup-aws-lambda/pkg/handler"
)

type application struct {
	Handler *handler.LambdaHandler
}

func (a *application) RequestHandler() interface{} {
	return a.Handler.HandleRequest
}

func SetupApp() *application {

	errors.SetDefaultErrorCode("999")

	// Initialize local clients, processor and lambda handler
	var (
		// awsSess        = aws.NewSession()
		// db             = dynamodb.New(awsSess)
		processorDummy = processor.NewProcessorDummy()
		lambdaHandler  = handler.NewLambdaHandler(processorDummy)
	)

	return &application{
		Handler: lambdaHandler,
	}
}
