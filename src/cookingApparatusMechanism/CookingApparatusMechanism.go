package cookingApparatusMechanism

import (
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/icook"
	"sync"
)

// region Singleton

var mechs Container

func Set(array Container) {
	mechs = array
}

func Get() Container {
	return mechs
}

// endregion

type CookingApparatusMechanism struct {
	Apparatus string

	busy bool

	queue  []icook.ICook
	locker chan bool

	lock sync.Mutex
}

type Container []*CookingApparatusMechanism

func New(apparatus string) *CookingApparatusMechanism {
	return &CookingApparatusMechanism{
		Apparatus: apparatus,
		busy:      false,
		queue:     make([]icook.ICook, 0),
		locker:    make(chan bool),
	}
}

func (m *CookingApparatusMechanism) AddQueue(item icook.ICook) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.queue = append(m.queue, item)
}

func (m *CookingApparatusMechanism) LenQueue() int {
	m.lock.Lock()
	defer m.lock.Unlock()

	return len(m.queue)
}

func (m *CookingApparatusMechanism) getFirst() icook.ICook {
	if len(m.queue) == 0 {
		return nil
	}
	return m.queue[0]
}

func (m *CookingApparatusMechanism) GetStatus() bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.busy
}

func (m *CookingApparatusMechanism) SetBusy(item icook.ICook) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.busy {
		return false
	}

	if m.getFirst() != item {
		return false
	}

	m.busy = true

	m.queue = m.queue[1:]

	return true
}

func (m *CookingApparatusMechanism) SetEmpty() {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.busy {
		return
	}

	m.busy = false
}
