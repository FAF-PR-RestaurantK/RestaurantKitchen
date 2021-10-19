package cookThread

import (
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/configuration"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/cookingApparatusMechanism"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/icook"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/item"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/queue"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/utils"
	"log"
	"math"
	"time"
)

type CookThread struct {
	cook             icook.ICook
	queue            *ThreadCookDetailsHeap
	currentItem      *ThreadCookDetails
	currentMech      *cookingApparatusMechanism.CookingApparatusMechanism
	currentItemTimer <-chan time.Time

	queueLeftTime time.Duration
	stash         *queue.Queue

	status bool
}

func New(cook icook.ICook) *CookThread {
	return &CookThread{
		cook:          cook,
		queue:         &ThreadCookDetailsHeap{},
		stash:         queue.New(),
		currentItem:   nil,
		queueLeftTime: 0,
		status:        false,
	}
}

// region Public props

func (thread *CookThread) GetId() int {
	return thread.cook.GetId()
}

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

	cookDetail := ThreadCookDetails{
		item:   item,
		detail: detail,
	}

	thread.queue.Push(cookDetail)
}

func (thread *CookThread) Update() {
	if thread.currentItem == nil {
		thread.popItem()
		return
	}

	if thread.status == false {
		thread.changeStatus()
		return
	}

	select {
	case <-thread.currentItemTimer:
		thread.currentItem.detail.CookID = thread.cook.GetId()
		thread.currentItem.detail.ReadyStatus = true
		thread.status = false
		thread.freeMechanism()
		thread.popItem()
		return
	default:
		return
	}
}

// endregion

// region Private methods

func (thread *CookThread) popItem() {
	if thread.queue.Len() == 0 && thread.stash.Len() == 0 {
		thread.currentItem = nil
		thread.currentMech = nil
		return
	}

	thread.currentMech = nil

	cookDetail := thread.queue.Pop().(ThreadCookDetails)
	thread.currentItem = &cookDetail
	thread.setApparatus(&cookDetail)

	thread.status = false
}

func (thread *CookThread) changeStatus() {
	if thread.currentMech != nil {
		if thread.currentMech.GetStatus() {
			return
		}

		if !thread.currentMech.SetBusy(thread) {
			return
		}
	}

	itemElem := thread.currentItem.item
	thread.queueLeftTime -= itemElem.Duration
	thread.currentItemTimer = time.After(itemElem.Duration)

	thread.status = true
}

func (thread *CookThread) setApparatus(detail *ThreadCookDetails) {
	if detail.item.CookingApparatus == nil {
		return
	}

	thread.currentMech = nil

	var currentMech *cookingApparatusMechanism.CookingApparatusMechanism = nil
	var minLen int = math.MaxInt64

	mechs := cookingApparatusMechanism.Get()
	for _, mech := range mechs {
		if *detail.item.CookingApparatus == mech.Apparatus {
			lenQ := mech.LenQueue()
			if lenQ < minLen {
				currentMech = mech
				minLen = lenQ
			}
		}
	}

	if currentMech == nil {
		log.Fatalf("exit: %s\n", "Unknown `cooking apparatus` in calculation.")
	}

	currentMech.AddQueue(thread)

	thread.currentMech = currentMech
}

func (thread *CookThread) freeMechanism() {
	if thread.currentMech != nil {
		thread.currentMech.SetEmpty()
	}
}

// endregion
