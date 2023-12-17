package pq

// package pq provides a non-intrusive priority queue.  Unlike other code in
// ick, something like this should be in the standard library (and is lacking
// because generics are so new).

import (
	"container/heap"
)

// Items in a PriorityQueue must implement the Priority interface.
type Prioritizer interface {
	Priority() int
}

// pqItem represents an item, and its position, in a priority queue.  This
// struct makes this a non-intrusive priority queue, different than what is
// described at https://pkg.go.dev/container/heap#example-package-PriorityQueue
type pqItem[T Prioritizer] struct {
	value T
	index int
}

// pqImpl provides a type to hang our receivers on.  This implements all
// methods in heap.Interface.
type pqImpl[T Prioritizer] []*pqItem[T]

func (q *pqImpl[T]) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item.value
}

func (q *pqImpl[T]) Push(x any) {
	n := len(*q)
	item := x.(pqItem[T])
	item.index = n
	*q = append(*q, &item)
}

func (pq pqImpl[T]) Len() int { return len(pq) }

// Return items with the lowest priority value.
func (pq pqImpl[T]) Less(i, j int) bool {
	return pq[i].value.Priority() < pq[j].value.Priority()
}

func (pq pqImpl[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// PriorityQueue provides a priority queue based on container/heap.
//
// An update method is missing because doing that requires access to
// the item to fix the heap after an item changes value.
type PriorityQueue[T Prioritizer] struct {
	pq *pqImpl[T]
}

func New[T Prioritizer]() *PriorityQueue[T] {
	n := &PriorityQueue[T]{pq: &pqImpl[T]{}}
	heap.Init(n.pq)
	return n
}

func (q *PriorityQueue[T]) Push(t T) {
	item := pqItem[T]{value: t}
	heap.Push(q.pq, item)
}

func (q *PriorityQueue[T]) Pop() T {
	return heap.Pop(q.pq).(T)
}

func (q *PriorityQueue[T]) Len() int {
	return q.pq.Len()
}
