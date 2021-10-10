package utils

import "time"

type DistributionData struct {
	OrderID        int              `json:"order-id"`
	TableID        int              `json:"table-id"`
	WaiterID       int              `json:"waiter-id"`
	Items          []int            `json:"items"`
	Priority       int              `json:"priority"`
	MaxWait        float32          `json:"max-wait"`
	PickUpTime     int64            `json:"pick-up-time"`
	CookingTime    int              `json:"cooking-time"`
	CookingDetails []CookingDetails `json:"cooking-details"`

	receivedTime time.Time
}

func NewDistData(order *OrderData) *DistributionData {
	return &DistributionData{
		OrderID:        order.OrderID,
		TableID:        order.TableID,
		WaiterID:       order.WaiterID,
		Items:          order.Items,
		Priority:       order.Priority,
		MaxWait:        order.MaxWait,
		PickUpTime:     order.PickUpTime,
		CookingDetails: make([]CookingDetails, 0, len(order.Items)),
	}
}

func (data *DistributionData) SetReceivedTime(value time.Time) {
	data.receivedTime = value
}

func (data *DistributionData) GetReceivedTime() time.Time {
	return data.receivedTime
}
