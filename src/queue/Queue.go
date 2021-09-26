package queue

import "container/list"

type Queue struct {
	list *list.List
}

func New() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (q *Queue) Push(e interface{}) {
	q.list.PushBack(e)
}

func (q *Queue) Pop() interface{} {
	if q.list.Len() <= 0 {
		return nil
	}

	e := q.list.Front()
	defer q.list.Remove(e)

	return e.Value
}

func (q *Queue) Len() int {
	return q.list.Len()
}
