package cookThread

import (
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/configuration"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/icook"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/item"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/queue"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/utils"
	"time"
)

type cookDetails struct {
	item   *item.Item
	detail *utils.CookingDetails
}

type CookThread struct {
	cook             icook.ICook
	queue            *queue.Queue
	currentItem      *cookDetails
	currentItemTimer <-chan time.Time

	queueLeftTime time.Duration
}

func New(cook icook.ICook) *CookThread {
	return &CookThread{
		cook:          cook,
		queue:         queue.New(),
		currentItem:   nil,
		queueLeftTime: 0,
	}
}

// region Public props

func (thread *CookThread) GetProficiency() int {
	return thread.cook.GetProficiency()
}

func (thread *CookThread) GetTimeLeft() time.Duration {
	if thread.currentItem == nil {
		return thread.queueLeftTime
	}
	return thread.queueLeftTime + thread.currentItem.item.Duration
}

// endregion

// region Public methods

func (thread *CookThread) PushItem(item *item.Item, detail *utils.CookingDetails) {
	item.Duration = time.Duration(item.PreparationTime) * configuration.TimeUnit
	thread.queueLeftTime += item.Duration

	cookDetail := cookDetails{
		item:   item,
		detail: detail,
	}

	if thread.queue == nil {
		thread.queue = queue.New()
	}

	thread.queue.Push(cookDetail)
}

func (thread *CookThread) Update() {
	if thread.currentItem == nil {
		thread.popItem()
		return
	}

	select {
	case <-thread.currentItemTimer:
		thread.currentItem.detail.CookID = thread.cook.GetId()
		thread.currentItem.detail.ReadyStatus = true
		thread.popItem()
		return
	default:
		return
	}
}

// endregion

// region Private methods

func (thread *CookThread) popItem() {
	if thread.queue.Len() == 0 {
		thread.currentItem = nil
		return
	}

	cookDetail := thread.queue.Pop().(cookDetails)

	itemElem := cookDetail.item

	itemDuration := time.Duration(itemElem.PreparationTime) * configuration.TimeUnit
	thread.queueLeftTime -= itemDuration

	thread.currentItemTimer = time.After(itemDuration)
	thread.currentItem = &cookDetail
}

// endregion
