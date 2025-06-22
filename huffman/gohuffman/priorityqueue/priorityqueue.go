package priorityqueue

import "container/heap"

type Item interface {
	Priority() int
}

type Queue[T Item] struct {
	items []T
}

func New[T Item]() *Queue[T] {
	return &Queue[T]{items: []T{}}
}

func (pq Queue[T]) Len() int {
	return len(pq.items)
}

func (pq Queue[T]) Less(i, j int) bool {
	return pq.items[i].Priority() < pq.items[j].Priority()
}

func (pq Queue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *Queue[T]) Push(x any) {
	pq.items = append(pq.items, x.(T))
}

func (pq *Queue[T]) Pop() any {
	n := len(pq.items)
	item := pq.items[n-1]
	pq.items = pq.items[:n-1]
	return item
}

func (pq *Queue[T]) Enqueue(item T) {
	heap.Push(pq, item)
}

func (pq *Queue[T]) Dequeue() T {
	return heap.Pop(pq).(T)
}

func (pq *Queue[T]) Peek() T {
	return pq.items[0]
}

func (pq *Queue[T]) Size() int {
	return len(pq.items)
}

func (pq *Queue[T]) Init() {
	heap.Init(pq)
}