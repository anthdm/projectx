package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {
	l := NewList[int]()
	assert.Equal(t, l.Data, []int{})
}

func TestListClear(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
	}

	assert.Equal(t, n, l.Len())
	l.Clear()
	assert.Equal(t, 0, l.Len())
}

func TestListContains(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
		assert.True(t, l.Contains(i))
	}
}

func TestListGetIndex(t *testing.T) {
	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)
		assert.Equal(t, l.GetIndex(data), i)
	}

	assert.Equal(t, l.GetIndex("bar"), -1)
}

func TestListRemove(t *testing.T) {
	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)
		l.Remove(data)
		assert.Equal(t, l.Contains(data), false)
	}

	assert.Equal(t, l.Len(), 0)
}

func TestListGet(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
		assert.True(t, l.Contains(i))
		assert.Equal(t, l.Get(i), i)
	}
}

func TestRemoveAt(t *testing.T) {
	l := NewList[int]()
	l.Insert(1)
	l.Insert(2)
	l.Insert(3)
	l.Insert(4)

	l.Pop(0)
	assert.Equal(t, l.Get(0), 2)
}

func TestListAdd(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
	}

	assert.Equal(t, n, l.Len())
}

func TestListLast(t *testing.T) {
	l := NewList[int]()
	l.Insert(1)
	l.Insert(2)
	l.Insert(3)

	assert.Equal(t, 3, l.Last())
}
