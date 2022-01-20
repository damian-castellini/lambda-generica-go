package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	app "users-backup-aws-lambda/internal"
)

func main() {
	application := app.SetupApp()
	lambda.Start(application.RequestHandler())
}
