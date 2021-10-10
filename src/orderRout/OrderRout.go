package orderRout

import (
	"encoding/json"
	"fmt"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/orderManager"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/utils"
	"io"
	"net/http"
)

func OrderHandler(writer http.ResponseWriter, request *http.Request) {
	var data utils.OrderData
	var response string

	jsonData, err := io.ReadAll(request.Body)
	str := string(jsonData)
	fmt.Println(str)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	http.Error(writer, response, http.StatusOK)

	orderManager.PushOrder(utils.NewDistData(&data))
}
