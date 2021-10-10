package item

import "time"

type Item struct {
	Id               int              `json:"id"`
	Name             string           `json:"name"`
	PreparationTime  int              `json:"preparation-time"`
	Complexity       int              `json:"complexity"`
	CookingApparatus CookingApparatus `json:"cooking-apparatus"`

	Duration time.Duration
	Priority int
}

func GetItem(id int, items []Item) *Item {
	for i := range items {
		if items[i].Id == id {
			return &items[i]
		}
	}

	return nil
}
