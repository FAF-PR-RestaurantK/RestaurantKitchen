package sendRequest

import (
	"encoding/json"
	"fmt"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/configuration"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/utils"
)

const (
	HttpAddr = "http://"
)

func SendDistribution(order *utils.DistributionData, conf *configuration.Configuration) {
	addr := HttpAddr + conf.DinnerHallAddr + conf.DistributionRout

	fmt.Print("send: ")
	fmt.Println(order)

	jsonBuff, err := json.Marshal(*order)
	if err != nil {
		fmt.Println(err)
		return
	}

	SendRequest(addr, jsonBuff)
}
