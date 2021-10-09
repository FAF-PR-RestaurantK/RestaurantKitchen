package cook

import (
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/cookThread"
	"time"
)

type Cook struct {
	Rank        int
	Proficiency int
	Name        string
	CatchPhrase string
	Id          int

	threads []*cookThread.CookThread

	lastTime int64
}

func New(
	rank int,
	prof int,
	name string,
	phrase string,
	threadsCount int) *Cook {

	cook := &Cook{
		Rank:        rank,
		Proficiency: prof,
		Name:        name,
		CatchPhrase: phrase,

		threads: make([]*cookThread.CookThread, threadsCount),
	}

	for i := range cook.threads {
		cook.threads[i] = cookThread.New(cook)
	}

	return cook
}

// region Public property

func (cook *Cook) GetRank() int {
	return cook.Rank
}

func (cook *Cook) GetProficiency() int {
	return cook.Proficiency
}

func (cook *Cook) GetName() string {
	return cook.Name
}

func (cook *Cook) GetCatchPhrase() string {
	return cook.CatchPhrase
}

func (cook *Cook) GetThreads() []*cookThread.CookThread {
	return cook.threads
}

func (cook *Cook) SetThreads(count int) {
	if cook.threads == nil {
		cook.threads = make([]*cookThread.CookThread, count)

		for i := range cook.threads {
			cook.threads[i] = cookThread.New(cook)
		}
	}
}

func (cook *Cook) GetId() int {
	return cook.Id
}

// endregion

// region Public methods

func (cook *Cook) Run() {
	cook.lastTime = time.Now().Unix()

	for {
		cook.update()
	}
}

func (cook *Cook) update() {
	timeNow := time.Now().Unix()

	cook.lastTime = timeNow

	for _, thread := range cook.threads {
		thread.Update()
	}
}

// endregion
