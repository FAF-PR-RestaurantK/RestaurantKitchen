package main

import (
	"log"
	"net/http"
	"restaurant_client/src/handler"
)

func main() {
	http.HandleFunc("/distribution", handler.DistributionHandler)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}