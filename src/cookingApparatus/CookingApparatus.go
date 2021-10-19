package cookingApparatus

import "github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/item"

type CookingApparatus struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type Container map[string]int

func (c Container) Check(items *item.Container) bool {
	for _, i := range items.GetList() {
		apparatusRef := i.CookingApparatus
		if apparatusRef == nil {
			continue
		}
		_, ok := c[*apparatusRef]

		if ok == false {
			return false
		}
	}

	return true
}
