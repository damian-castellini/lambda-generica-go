package internal

import (
	"github.com/Bancar/uala-go-platform-product-dependencies/pkg/errors"
	_ "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"users-backup-aws-lambda/internal/aws"
	_ "users-backup-aws-lambda/internal/aws"
	"users-backup-aws-lambda/internal/processor"
	"users-backup-aws-lambda/internal/secret"
	"users-backup-aws-lambda/pkg/handler"
	"users-backup-aws-lambda/storage"
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
		awsSess            = aws.NewSession()
		svc                = s3.New(awsSess)
		dbSecret           = secret.NewSecret()
		databaseConnection = storage.NewDatabaseConnection()
		processor          = processor.NewProcessor(dbSecret, svc, databaseConnection)
		lambdaHandler      = handler.NewLambdaHandler(processor)
	)

	return &application{
		Handler: lambdaHandler,
	}
}
