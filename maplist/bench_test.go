package maplist

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

// 证明 maplist add/del 操作时间复杂度是 O(1)
func BenchmarkMapListS(b *testing.B) {
	for _, size := range []int{5, 10, 100, 200, 1000, 10000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			_benchMapListS_WithSize(b, size)
		})
	}
}

func _benchMapListS_WithSize(b *testing.B, size int) {
	// build := func(size int) *MapListS {
	// 	list := []string{}
	// 	mapList := NewS(&list)
	// 	for i := 0; i < size; i++ {
	// 		mapList.Add(strconv.Itoa(i))
	// 	}
	// 	return mapList
	// }

	build_ids := func(size int) []string {
		list := make([]string, size)
		for i := 0; i < size; i++ {
			list[i] = strconv.Itoa(i)
		}
		return list
	}

	case_del := func(src *[]string) *MapList[string] {
		ids := (*src)[:]
		obj := New(src)
		for _, id := range ids {
			if !obj.Remove(id) {
				b.Logf("debug")
			}
		}
		return obj
	}

	case_add := func(src []string) *MapList[string] {
		l := []string{}
		obj := New(&l)
		for i := 0; i < len(src); i++ {
			obj.Add(src[i])
		}
		return obj
	}

	ids := build_ids(size)
	b.Run(fmt.Sprintf("添加%d个元素", size), func(b *testing.B) {
		src := ids[:]
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			obj := case_add(src)
			if obj.Length() != size {
				b.FailNow()
			}
		}
		b.ReportMetric(float64(size*b.N), "总次数/秒")
	})
	shuffle := func(src []string) {
		rand.Shuffle(len(src), func(i, j int) {
			src[i], src[j] = src[j], src[i]
		})
	}
	clone := func(src []string) []string {
		cloned := make([]string, len(src))
		if copy(cloned, src) != len(src) {
			panic(fmt.Errorf("copy failed"))
		}
		return cloned
	}

	ids = build_ids(size)
	b.Run(fmt.Sprintf("删除%d个元素(含clone+shuffle)", size), func(b *testing.B) {
		// b.ResetTimer()
		for i := 0; i < b.N; i++ {
			src := clone(ids)
			shuffle(src)
			obj := case_del(&src)
			if obj.Length() != 0 {
				b.Logf("i:%d obj.length:%d src.length:%d err:%v", i, obj.Length(), len(src), obj.CheckSlowly())
				b.FailNow()
			}
			if len(src) != 0 {
				b.Logf("err-2: obj.length:%d src.length:%d", obj.Length(), len(src))
				b.FailNow()
			}
		}
		b.ReportMetric(float64(size*b.N), "总次数/秒")
	})

}

func Benchmark_map(b *testing.B) {
	v := new(int)
	*v = 1
	gen := func(size int) map[int64]*int {
		m := map[int64]*int{}
		for i := 0; i < size; i++ {
			m[int64(i)] = v
		}
		return m
	}
	for _, v := range []int{10_000, 100_000, 1_000_000} {
		m := gen(v)
		b.Run(fmt.Sprintf("size=%d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v := m[int64(i)]
				_ = v
			}
		})
	}

}
