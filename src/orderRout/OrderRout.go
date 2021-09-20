package orderRout

import (
	"encoding/json"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/utils"
	"io"
	"net/http"
)

func OrderHandler(writer http.ResponseWriter, request *http.Request) {
	var data utils.OrderData
	var response string

	jsonData, err := io.ReadAll(request.Body)
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
}
