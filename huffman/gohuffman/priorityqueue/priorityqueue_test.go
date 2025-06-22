package priorityqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestItem struct {
	value    int
	priority int
}

func (t TestItem) Priority() int {
	return t.priority
}

func TestPriorityQueue(t *testing.T) {
	pq := New[TestItem]()
	
	pq.Enqueue(TestItem{value: 1, priority: 3})
	pq.Enqueue(TestItem{value: 2, priority: 1})
	pq.Enqueue(TestItem{value: 3, priority: 2})
	
	assert.Equal(t, 3, pq.Size())
	
	first := pq.Dequeue()
	assert.Equal(t, 1, first.priority)
	assert.Equal(t, 2, first.value)
	
	second := pq.Dequeue()
	assert.Equal(t, 2, second.priority)
	assert.Equal(t, 3, second.value)
	
	third := pq.Dequeue()
	assert.Equal(t, 3, third.priority)
	assert.Equal(t, 1, third.value)
	
	assert.Equal(t, 0, pq.Size())
}