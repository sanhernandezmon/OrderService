package mappers

import (
	"OrderService/domain"
	"github.com/google/uuid"
)

func MapRequestToOrder(request domain.CreateOrderRequest) domain.Order {
	newUUID := uuid.New().String()
	return domain.Order{
		newUUID,
		request.UserID,
		request.Item,
		request.Quantity,
		request.TotalPrice}
}
func MapUUIDandPriceIntoEvent(orderId string, totalPrice int64) domain.CreateOrderEvent {
	return domain.CreateOrderEvent{orderId, totalPrice}
}
