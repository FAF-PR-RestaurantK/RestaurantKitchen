package cookThread

import (
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/item"
	"github.com/FAF-PR-RestaurantK/RestaurantKitchen/src/utils"
)

type ThreadCookDetails struct {
	item   *item.Item
	detail *utils.CookingDetails
}

type ThreadCookDetailsHeap []ThreadCookDetails

func (heap ThreadCookDetailsHeap) Len() int {
	return len(heap)
}

func (heap ThreadCookDetailsHeap) Less(i, j int) bool {
	return heap[i].item.Priority > heap[j].item.Priority
}

func (heap ThreadCookDetailsHeap) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}

func (heap *ThreadCookDetailsHeap) Push(x interface{}) {
	*heap = append(*heap, x.(ThreadCookDetails))
}

func (heap *ThreadCookDetailsHeap) Pop() interface{} {
	old := *heap
	count := len(old)

	x := old[count-1]
	*heap = old[0 : count-1]

	return x
}
