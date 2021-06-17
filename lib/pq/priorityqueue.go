// Package pq implements a priority queue data structure on top of container/heap.
// As an addition to regular operations, it allows an update of an items priority,
// allowing the queue to be used in graph search algorithms like Dijkstra's algorithm.
// Computational complexities of operations are mainly determined by container/heap.
// In addition, a map of items is maintained, allowing O(1) lookup needed for priority updates,
// which themselves are O(log n).
package pq

import (
	"container/heap"
	"errors"
)

// PriorityQueue represents the queue
type PriorityQueue struct {
	itemHeap *itemHeap
	lookup   map[Queuable]*item
	maxSize  int
}

type Queuable interface {
	Clean()
}

// New initializes an empty priority queue.
func New() PriorityQueue {
	return PriorityQueue{
		itemHeap: &itemHeap{},
		lookup:   make(map[Queuable]*item),
		maxSize:  -1,
	}
}

func NewWithMaxSize(maxSize int) PriorityQueue {
	return PriorityQueue{
		itemHeap: &itemHeap{},
		lookup:   make(map[Queuable]*item),
		maxSize:  maxSize,
	}
}

// Len returns the number of elements in the queue.
func (p *PriorityQueue) Len() int {
	return p.itemHeap.Len()
}

// Insert inserts a new element into the queue. No action is performed on duplicate elements.
func (p *PriorityQueue) Insert(v Queuable, priority, secondaryPrio float64) {
	_, ok := p.lookup[v]
	if ok {
		return
	}

	newItem := &item{
		value:         v,
		priority:      priority,
		secondaryPrio: secondaryPrio,
	}

	heap.Push(p.itemHeap, newItem)
	p.lookup[v] = newItem

	if p.maxSize > 0 && p.Len() > p.maxSize {
		old, _ := p.popEnd()

		if old != nil {
			old.Clean()
		}
	}
}

// Pop removes the element with the highest priority from the queue and returns it.
// In case of an empty queue, an error is returned.
func (p *PriorityQueue) Pop() (Queuable, error) {
	if len(*p.itemHeap) == 0 {
		return nil, errors.New("empty queue")
	}

	item := heap.Pop(p.itemHeap).(*item)
	delete(p.lookup, item.value)
	return item.value, nil
}

func (p *PriorityQueue) popEnd() (Queuable, error) {
	if len(*p.itemHeap) == 0 {
		return nil, errors.New("empty queue")
	}

	item := heap.Remove(p.itemHeap, p.Len()-1).(*item)
	delete(p.lookup, item.value)
	return item.value, nil
}

// UpdatePriority changes the priority of a given item.
// If the specified item is not present in the queue, no action is performed.
func (p *PriorityQueue) UpdatePriority(x Queuable, newPriority float64) {
	item, ok := p.lookup[x]
	if !ok {
		return
	}

	item.priority = newPriority
	heap.Fix(p.itemHeap, item.index)
}

type itemHeap []*item

type item struct {
	value         Queuable
	priority      float64
	secondaryPrio float64
	index         int
}

func (ih *itemHeap) Len() int {
	return len(*ih)
}

func (ih *itemHeap) Less(i, j int) bool {
	if (*ih)[i].priority == (*ih)[j].priority {
		return (*ih)[i].secondaryPrio < (*ih)[j].secondaryPrio
	}

	return (*ih)[i].priority < (*ih)[j].priority
}

func (ih *itemHeap) Swap(i, j int) {
	(*ih)[i], (*ih)[j] = (*ih)[j], (*ih)[i]
	(*ih)[i].index = i
	(*ih)[j].index = j
}

func (ih *itemHeap) Push(x interface{}) {
	it := x.(*item)
	it.index = len(*ih)
	*ih = append(*ih, it)
}

func (ih *itemHeap) Pop() interface{} {
	old := *ih
	item := old[len(old)-1]
	*ih = old[0 : len(old)-1]
	return item
}
