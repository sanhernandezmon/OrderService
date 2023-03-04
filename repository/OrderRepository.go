package repository

import (
	"OrderService/domain"
	"OrderService/mappers"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

func SaveOrderToDynamoDB(request domain.CreateOrderRequest) (string, error) {
	var order = mappers.MapRequestToOrder(request)
	err := AddElement(order)
	if err != nil {
		panic(err)
		return "", err
	}
	return order.OrderId, err
}

func SendOrderSQSMessage(orderID string, totalPrice int64) {
	var orderEvent = mappers.MapUUIDandPriceIntoEvent(orderID, totalPrice)
	message, err := json.Marshal(orderEvent)
	if err != nil {
		panic(err)
	}
	queueURL := "http://localhost:9324/queue/default"
	sqsURL := "http://localhost:9324"
	sqsClient := newSQS(endpoints.UsEast1RegionID, sqsURL)
	print("sending message to sqs")
	sendMessage(sqsClient, string(message), queueURL)
}
