package main

import (
	"encoding/json"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/configuration"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/cook"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/cookingApparatus"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/cookingApparatusMechanism"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/item"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/orderManager"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/orderRout"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	ConfPath             = "./conf/configuration.json"
	ItemsPath            = "./conf/items.json"
	CooksPath            = "./conf/cooks.json"
	CookingApparatusPath = "./conf/cookingApparatus.json"
)

func main() {
	conf := GetConf()
	container := GetItemContainer()
	cooks := GetCooks()
	cookingApparatuses := GetCookingApparatus()
	mechanisms := GetCookingApparatusMechanism(cookingApparatuses)

	ok := cookingApparatuses.Check(container)
	if ok == false {
		log.Fatalf("exit: %s\n", "Item container uses unknown `cooking apparatus`")
		return
	}

	cookingApparatusMechanism.Set(mechanisms)

	timeUnitMillisecondMultiplier := time.Duration(conf.TimeUnitMillisecondMultiplier)
	configuration.TimeUnit = time.Millisecond * timeUnitMillisecondMultiplier

	for i := range cooks {
		cooks[i].Id = i
		cooks[i].SetThreads(cooks[i].Proficiency)
	}

	orderManager.SetItems(container.GetList())
	orderManager.SetConf(&conf)
	orderManager.SetCooks(cooks)

	go orderManager.Get().Run()

	for i := range cooks {
		go cooks[i].Run()
	}

	http.HandleFunc(conf.OrderRout, orderRout.OrderHandler)

	err := http.ListenAndServe(conf.KitchenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConf() configuration.Configuration {
	var conf configuration.Configuration

	confFile, _ := os.Open(ConfPath)
	defer func(confFile *os.File) {
		_ = confFile.Close()
	}(confFile)

	jsonData, err := io.ReadAll(confFile)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return conf
	}

	err = json.Unmarshal(jsonData, &conf)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return conf
	}

	return conf
}

func GetItemContainer() *item.Container {
	var itemList []item.Item

	itemListFile, _ := os.Open(ItemsPath)
	defer func(itemListFile *os.File) {
		_ = itemListFile.Close()
	}(itemListFile)

	jsonData, err := io.ReadAll(itemListFile)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	err = json.Unmarshal(jsonData, &itemList)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	return item.NewContainer(itemList)
}

func GetCooks() []*cook.Cook {
	var cooks []*cook.Cook

	cooksFile, _ := os.Open(CooksPath)
	defer func(itemListFile *os.File) {
		_ = itemListFile.Close()
	}(cooksFile)

	jsonData, err := io.ReadAll(cooksFile)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	err = json.Unmarshal(jsonData, &cooks)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	return cooks
}

func GetCookingApparatus() cookingApparatus.Container {
	var apparatusArray []cookingApparatus.CookingApparatus

	cookingApparatusFile, _ := os.Open(CookingApparatusPath)
	defer func(itemListFile *os.File) {
		_ = itemListFile.Close()
	}(cookingApparatusFile)

	jsonData, err := io.ReadAll(cookingApparatusFile)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	err = json.Unmarshal(jsonData, &apparatusArray)
	if err != nil {
		log.Fatalf("exit: %s\n", err.Error())
		return nil
	}

	apparatusMap := make(cookingApparatus.Container)

	for i := range apparatusArray {
		apparatusMap[apparatusArray[i].Key] = apparatusArray[i].Value
	}

	return apparatusMap
}

func GetCookingApparatusMechanism(source cookingApparatus.Container) cookingApparatusMechanism.Container {
	container := make(cookingApparatusMechanism.Container, 0)

	for apparatus := range source {
		count := source[apparatus]
		for i := 0; i < count; i++ {
			mechanism := cookingApparatusMechanism.New(apparatus)
			container = append(container, mechanism)
		}
	}

	return container
}
