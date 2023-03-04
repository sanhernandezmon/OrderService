package repository

import (
	"OrderService/domain"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

// CreateLocalClient Creates a local DynamoDb Client on the specified port. Useful for connecting to DynamoDB Local or
// LocalStack.

func CreateDynamoClient() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String("eu-central-1")},
	)
	if err != nil {
		// Handle Session creation error
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	return svc
}

func AddElement(order domain.Order) error {
	svc := CreateDynamoClient()
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Orders"),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	return err
}

func CreateTable() {
	tableName := "Orders"
	svc := CreateDynamoClient()

	attrDef := []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("order_id"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("user_id"),
			AttributeType: aws.String("S"),
		},
	}

	// Define the key schema for the table
	keySchema := []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("order_id"),
			KeyType:       aws.String("HASH"),
		},
		{
			AttributeName: aws.String("user_id"),
			KeyType:       aws.String("RANGE"),
		},
	}

	// Create the table input struct
	tableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions:  attrDef,
		KeySchema:             keySchema,
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(5), WriteCapacityUnits: aws.Int64(5)},
		TableName:             aws.String(tableName),
	}
	// Create the table
	_, err := svc.CreateTable(tableInput)
	if err != nil {
		panic(err)
	}

	// Wait for the table to become active
	err = svc.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
	if err != nil {
		panic(err)
	}

	fmt.Println("Created the table", tableName)
}
