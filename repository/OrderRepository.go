package repository

import (
	"OrderService/domain"
	"OrderService/mappers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SaveOrderToDynamoDB(request domain.CreateOrderRequest) (string, error) {
	var cfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	var dynamoClient = dynamodb.NewFromConfig(cfg)

	var order = mappers.MapRequestToOrder(request)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("orders"),
		Item: map[string]types.AttributeValue{
			"order_id": &types.AttributeValueMemberS{
				Value: order.OrderId,
			},
			"user_id": &types.AttributeValueMemberS{
				Value: order.UserID,
			},
			"item": &types.AttributeValueMemberS{
				Value: order.Item,
			},
			"quantity": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", order.Quantity),
			},
			"total_price": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%.2f", order.TotalPrice),
			},
		},
	}

	if _, err := dynamoClient.PutItem(context.TODO(), input); err != nil {
		return "", fmt.Errorf("failed to save order to DynamoDB: %w", err)
	}

	return order.OrderId, nil
}

func SendOrderSQSMessage(orderID string, totalPrice int64) error {
	var cfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	var orderEvent = mappers.MapUUIDandPriceIntoEvent(orderID, totalPrice)
	message, err := json.Marshal(orderEvent)
	if err != nil {
		panic(err)
	}
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue"
	input := &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: aws.String(string(message)),
	}
	_, err = sqs.NewFromConfig(cfg).SendMessage(context.Background(), input)
	return err
}
