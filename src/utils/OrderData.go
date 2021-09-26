package utils

type OrderData struct {
	OrderID    int     `json:"order-id"`
	TableID    int     `json:"table-id"`
	WaiterID   int     `json:"waiter-id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float32 `json:"max-wait"`
	PickUpTime int64   `json:"pick-up-time"`
}
