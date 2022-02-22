package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"users-backup-aws-lambda/pkg/dto"
)

type DynamoFailedMigrate struct {
	db *dynamodb.DynamoDB
}

func NewDynamoFailedMigrate(db *dynamodb.DynamoDB) *DynamoFailedMigrate {
	return &DynamoFailedMigrate{db: db}
}

func (r *DynamoFailedMigrate) InsertDynamo(user dto.User) error {
	userItem, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("Error marshalling user: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      userItem,
		TableName: aws.String(dto.DYNAMO_TABLE),
	}

	if _, err = r.db.PutItem(input); err != nil {
		return fmt.Errorf("Error inserting user: %w", err)
	}

	return nil
}
