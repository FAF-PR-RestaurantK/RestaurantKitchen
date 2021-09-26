package utils

type CookingDetails struct {
	FoodID int `json:"food_id"`
	CookID int `json:"cook_id"`

	ReadyStatus bool
}
