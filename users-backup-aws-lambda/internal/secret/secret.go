package secret

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"users-backup-aws-lambda/pkg/dto"
)

type (
	dbSecret struct{}
)

func NewSecret() *dbSecret {
	return &dbSecret{}
}

func (*dbSecret) GetDBSecret() string {
	sessionCreated, err := session.NewSession()
	if err != nil {
		fmt.Println(err.Error())
		return "Error creating session on GetDBSecret: " + err.Error()
	}
	svc := secretsmanager.New(sessionCreated)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(dto.SECRET_NAME),
	}
	fmt.Println(input)
	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "Error getting secret value on GetSecretValue: " + err.Error()
	}

	if result.SecretString != nil {
		return *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return "Error encoding secret value: " + err.Error()
		}
		return string(decodedBinarySecretBytes[:len])
	}
}
