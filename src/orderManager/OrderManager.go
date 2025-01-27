package orderManager

import (
	"fmt"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/configuration"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/cook"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/cookThread"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/item"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/queue"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/sendRequest"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/utils"
	"math"
	"sync"
	"time"
)

type OrderManager struct {
	queue  *queue.Queue
	orders map[*utils.OrderData][]bool

	workingOrders []*utils.DistributionData

	cooks []*cookThread.CookThread
	items []item.Item
	conf  *configuration.Configuration

	buffer chan bool
	lock   sync.Mutex
}

// region Static methods

var instance *OrderManager

func newOrderManager() *OrderManager {
	return &OrderManager{
		queue:  queue.New(),
		orders: make(map[*utils.OrderData][]bool),
	}
}

func Get() *OrderManager {
	if instance == nil {
		instance = newOrderManager()
	}

	return instance
}

func SetCooks(cooks []*cook.Cook) {
	manager := Get()
	for i := range cooks {
		threads := cooks[i].GetThreads()
		for j := range threads {
			thread := threads[j]
			manager.cooks = append(manager.cooks, thread)
		}
	}
}

func SetItems(items []item.Item) {
	manager := Get()
	manager.items = items
}

func SetConf(conf *configuration.Configuration) {
	manager := Get()
	manager.conf = conf

	manager.workingOrders = make([]*utils.DistributionData, 0, conf.TableCount)
	manager.buffer = make(chan bool, conf.OrderListLen)
}

func PushOrder(order *utils.DistributionData) {
	manager := Get()

	manager.lock.Lock()

	manager.buffer <- false

	fmt.Print("get: ")
	fmt.Println(order)

	order.SetReceivedTime(time.Now())
	manager.Provide(order)

	manager.lock.Unlock()
}

// endregion

// region Public methods

func (manager *OrderManager) Run() {
	for {
		manager.update()
	}
}

func (manager *OrderManager) Provide(data *utils.DistributionData) {
	manager.workingOrders = append(manager.workingOrders, data)

	manager.setupCookingDetails(data)

	for i := range data.CookingDetails {
		manager.sendItemCook(&data.CookingDetails[i], data.Priority)
	}
}

// endregion

// region Private methods

func (manager *OrderManager) update() {
	manager.outputDataProvide()
}

func (manager *OrderManager) outputDataProvide() {
	if len(manager.workingOrders) != 0 {
		for i, data := range manager.workingOrders {
			count := 0
			for _, detail := range data.CookingDetails {
				if detail.ReadyStatus == true {
					count += 1
				}
			}

			if count == len(data.CookingDetails) {
				manager.getPreparedTime(data)
				sendRequest.SendDistribution(data, manager.conf)
				manager.workingOrders = manager.remove(manager.workingOrders, i)
				<-manager.buffer
				return
			}
		}
	}
}

func (manager *OrderManager) getPreparedTime(data *utils.DistributionData) {
	receivedTime := data.GetReceivedTime()
	preparedTime := time.Since(receivedTime)
	cookingTime := preparedTime / configuration.TimeUnit
	data.CookingTime = int(cookingTime)
}

func (manager *OrderManager) setupCookingDetails(data *utils.DistributionData) {
	for i := range data.Items {
		cookingDetails := &utils.CookingDetails{
			FoodID:      data.Items[i],
			ReadyStatus: false,
		}
		data.CookingDetails = append(data.CookingDetails, *cookingDetails)
	}
}

func (manager *OrderManager) sendItemCook(cookingDetails *utils.CookingDetails, priority int) {
	itemElem := item.GetItem(cookingDetails.FoodID, manager.items)
	cooker := manager.getCook(itemElem)
	itemElem.Priority = priority

	cooker.PushItem(itemElem, cookingDetails)

}

func (manager *OrderManager) getCook(item *item.Item) *cookThread.CookThread {
	var currentCookThread *cookThread.CookThread = nil
	var minTime time.Duration = math.MaxInt64

	for _, thread := range manager.cooks {
		if thread.GetProficiency() >= item.Complexity {
			timeLeft := thread.GetTimeLeft()
			if timeLeft < minTime {
				currentCookThread = thread
				minTime = timeLeft
			}
		}
	}

	return currentCookThread
}

func (manager *OrderManager) remove(slice []*utils.DistributionData, s int) []*utils.DistributionData {
	return append(slice[:s], slice[s+1:]...)
}

// endregion
