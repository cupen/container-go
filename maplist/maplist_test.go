package maplist

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapList_nil(t *testing.T) {
	assert := assert.New(t)
	var list []string = nil
	obj := New(&list)
	assert.False(obj.Has(""))
	v, flag := obj.GetByIndex(0)
	assert.Equal(0, obj.Length())
	assert.Equal("", v)
	assert.False(flag)
	assert.False(obj.RemoveByIndex(0))
	assert.False(obj.Remove(""))
}

func TestMapList(t *testing.T) {
	// for _, size := range []int{10, 100, 1000, 10000} {
	for _, size := range []int{10} {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			_testMapList_WithSize(t, size)
		})
	}
}

func _testMapList_WithSize(t *testing.T, size int) {
	list := []uint64{}
	SIZE := size
	t.Run("有序删除 by index", func(t *testing.T) {
		assert := require.New(t)
		mapList := New(&list)
		for i := 0; i < SIZE; i++ {
			assert.Equal(i, mapList.Length())
			mapList.Add(uint64(i))
			assert.Equal(i+1, mapList.Length())
			assert.NoError(mapList.CheckSlowly())
			assert.Equal(len(list), mapList.Length())
		}

		for i := SIZE - 1; i >= 0; i-- {
			assert.Equal(i+1, mapList.Length())
			assert.Truef(mapList.RemoveByIndex(i), "i=%d length=%d", i, mapList.Length())
			assert.Equalf(i, mapList.Length(), "i=%d length=%d", i, mapList.Length())
			assert.NoError(mapList.CheckSlowly())
			assert.Equal(len(list), mapList.Length())
		}
	})

	t.Run("无序删除 by index", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal(0, len(list))
		mapList := New(&list)

		for i := 0; i < SIZE; i++ {
			mapList.Add(uint64(i))
			assert.Equal(i+1, mapList.Length())
			assert.NoError(mapList.CheckSlowly())
			assert.Equal(len(list), mapList.Length())
		}

		for i := 0; i < SIZE; i++ {
			index := rand.Intn(mapList.Length())
			mapList.RemoveByIndex(index)

			assert.Equalf(SIZE-i-1, mapList.Length(), "i=%d length=%d", i, mapList.Length())
			assert.NoError(mapList.CheckSlowly())
			assert.Equal(len(list), mapList.Length())
		}
	})

	t.Run("无序删除 by ID", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal(0, len(list))
		mapList := New(&list)

		for i := 0; i < SIZE; i++ {
			mapList.Add(uint64(i))
			assert.Equal(i+1, mapList.Length())
			assert.NoError(mapList.CheckSlowly())
			assert.Equal(len(list), mapList.Length())
		}

		for i, index := range rand.Perm(SIZE) {
			id := uint64(index)
			mapList.Remove(id)
			assert.Equalf(SIZE-i-1, mapList.Length(), "id=%s length=%d", id, mapList.Length())
			assert.NoError(mapList.CheckSlowly())
			assert.Equal(len(list), mapList.Length())
		}
	})
}
