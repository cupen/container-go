package maplist

import "fmt"

type value interface {
	string | int | int32 | int64 | uint | uint32 | uint64
}

// 此数据结构仅用于为 map 提供可支持顺序访问的 list 接口
// 虽同时具备 map 和 list 的部分优势,但会使用双倍内存.
// 时间复杂度O(1),空间复杂度O(n) "两倍
// 注：list 暂不保证有序,删除元素时会把末尾元素前移填充。
type MapList[T value] struct {
	list  *[]T
	m     map[T]int
	empty T
}

func New[T value](list *[]T) *MapList[T] {
	obj := MapList[T]{
		list: list,
	}
	obj.Rebuild()
	return &obj
}

func (fm *MapList[T]) Add(elem T) bool {
	if _, isExist := fm.m[elem]; isExist {
		return false
	}
	*fm.list = append(*fm.list, elem)
	fm.m[elem] = len(*fm.list) - 1
	return true
}

// O(1)
func (fm *MapList[T]) Has(elem T) bool {
	_, isExist := fm.m[elem]
	return isExist
}

// O(1)
func (fm *MapList[T]) Remove(elem T) bool {
	index, isExist := fm.m[elem]
	if !isExist {
		return false
	}
	return fm.RemoveByIndex(index)
}

// O(1)
func (fm *MapList[T]) GetByIndex(index int) (T, bool) {
	list := *fm.list
	if index < 0 || index >= len(list) {
		return fm.empty, false
	}
	return list[index], true
}

// O(1)
// 注：删除元素时会把末尾元素前移填充。
func (fm *MapList[T]) RemoveByIndex(index int) bool {
	deletedElem, isExist := fm.GetByIndex(index)
	if !isExist {
		return false
	}
	list := *fm.list
	lastIndex := len(list) - 1
	if index != lastIndex {
		lastElem := list[lastIndex]
		list[index] = lastElem
		fm.m[lastElem] = index
	}
	// 同步更新
	delete(fm.m, deletedElem)
	*fm.list = list[:lastIndex]
	return true
}

// O(1)
func (fm *MapList[T]) Length() int {
	return len(*fm.list)
}

// O(1)
func (fm *MapList[T]) CheckQuickly() error {
	if len(*fm.list) != len(fm.m) {
		return fmt.Errorf("invalid length of list[%d] and map[%d]", len(*fm.list), len(fm.m))
	}
	return nil
}

// O(n) 仅用于临时排查错误
func (fm *MapList[T]) CheckSlowly() error {
	if len(*fm.list) != len(fm.m) {
		return fmt.Errorf("invalid length of list[%d] and map[%d]", len(*fm.list), len(fm.m))
	}
	for i, v := range *fm.list {
		if fm.m[v] != i {
			return fmt.Errorf("invalid index of list and map. index=%d value=%v", i, v)
		}
	}
	return nil
}

// 仅用于自动纠错
func (fm *MapList[T]) RebuildIf() bool {
	if len(fm.m) != len(*fm.list) {
		fm.Rebuild()
		return true
	}
	return false
}

// 仅用于临时排查错误
func (fm *MapList[T]) Rebuild() {
	fm.m = fm.buildMap(*fm.list)
}

func (fm *MapList[T]) buildMap(list []T) map[T]int {
	m := map[T]int{}
	for i, T := range list {
		m[T] = i
	}
	return m
}
