package util

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceContains(t *testing.T) {
	t.Parallel()
	res1 := SliceContains([]int{0, 1, 2, 3, 4, 5}, 5)
	res2 := SliceContains([]int{0, 1, 2, 3, 4, 5}, 6)

	require.True(t, res1)
	require.False(t, res2)
}

func TestSliceContainsOneBy(t *testing.T) {
	t.Parallel()

	type a struct {
		A int
		B int
	}

	a1 := []a{{
		A: 1,
		B: 1,
	}, {
		A: 2,
		B: 2,
	}, {
		A: 3,
		B: 3,
	}}

	res1 := SliceContainsOneBy(a1, func(item a) bool {
		return item.A == 1 && item.B == 2
	})
	res2 := SliceContainsOneBy(a1, func(item a) bool {
		return item.A == 1 && item.B == 1
	})

	a2 := []string{"a", "b", "c"}
	res3 := SliceContainsOneBy(a2, func(item string) bool {
		return item == ""
	})
	res4 := SliceContainsOneBy(a2, func(item string) bool {
		return item == "b"
	})

	require.False(t, res1)
	require.True(t, res2)
	require.False(t, res3)
	require.True(t, res4)
}

func TestSliceContainsAll(t *testing.T) {
	t.Parallel()

	res1 := SliceContainsAll([]int{1, 2, 3, 4, 5}, []int{1, 2})
	res2 := SliceContainsAll([]int{1, 2, 3, 4, 5}, []int{1, 6})
	res3 := SliceContainsAll([]int{1, 2, 3, 4, 5}, []int{0, 6})
	res4 := SliceContainsAll([]int{1, 2, 3, 4, 5}, []int{})

	require.True(t, res1)
	require.False(t, res2)
	require.False(t, res3)
	require.True(t, res4)
}

func TestSliceContainsAllBy(t *testing.T) {
	t.Parallel()

	res1 := SliceContainsAllBy([]int{1, 2, 3}, func(i int) bool {
		return i < 4
	})
	res2 := SliceContainsAllBy([]int{1, 2, 3}, func(i int) bool {
		return i < 3
	})
	res3 := SliceContainsAllBy([]int{1, 2, 3}, func(i int) bool {
		return i < 0
	})
	res4 := SliceContainsAllBy([]int{}, func(i int) bool {
		return i < 4
	})

	require.True(t, res1)
	require.False(t, res2)
	require.False(t, res3)
	require.True(t, res4)
}

func TestSliceContainsOneOf(t *testing.T) {
	t.Parallel()

	res1 := SliceContainsOneOf([]int{1, 2, 3}, []int{0, 4})
	res2 := SliceContainsOneOf([]int{1, 2, 3}, []int{1, 4})
	res3 := SliceContainsOneOf([]int{1, 2, 3}, []int{})
	res4 := SliceContainsOneOf([]int{1, 2, 3}, []int{1, 2, 3, 4})
	res5 := SliceContainsOneOf([]int{}, []int{1, 2, 3, 4})

	require.False(t, res1)
	require.True(t, res2)
	require.False(t, res3)
	require.True(t, res4)
	require.False(t, res5)

	type a struct {
		A int
	}

	a1 := &a{A: 1}
	res6 := SliceContainsOneOf([]*a{a1, {A: 2}}, []*a{a1})
	res7 := SliceContainsOneOf([]*a{a1, {A: 2}}, []*a{{A: 1}})
	require.True(t, res6)
	require.False(t, res7)
}

func TestSliceContainsNoneBy(t *testing.T) {
	t.Parallel()

	res1 := SliceContainsNoneBy([]int{1, 2, 3}, func(i int) bool {
		return i < 4
	})
	res2 := SliceContainsNoneBy([]int{1, 2, 3}, func(i int) bool {
		return i < 3
	})
	res3 := SliceContainsNoneBy([]int{1, 2, 3}, func(i int) bool {
		return i < 0
	})
	res4 := SliceContainsNoneBy([]int{}, func(i int) bool {
		return i < 4
	})

	require.False(t, res1)
	require.False(t, res2)
	require.True(t, res3)
	require.True(t, res4)
}

func TestSliceIntersect(t *testing.T) {
	t.Parallel()

	res1 := SliceIntersect([]int{1, 2, 3})
	res2 := SliceIntersect([]int{1, 2, 3}, []int{})
	res3 := SliceIntersect([]int{1, 2, 3}, []int{1, 2})
	res4 := SliceIntersect([]int{1, 2, 3}, []int{3, 4, 5})
	res5 := SliceIntersect([]int{}, []int{0, 1, 2})
	res6 := SliceIntersect([]int{1, 2, 1}, []int{1, 2, 3})
	res7 := SliceIntersect([]int{1, 2}, []int{2, 3}, []int{3, 4, 5}, []int{4, 5, 6})
	res8 := SliceIntersect([]int{1, 2}, []int{2, 3}, []int{3, 4, 2}, []int{4, 5, 2})

	require.Equal(t, []int{1, 2, 3}, res1)
	require.Equal(t, []int{}, res2)
	require.Equal(t, []int{1, 2}, res3)
	require.Equal(t, []int{3}, res4)
	require.Equal(t, []int{}, res5)
	require.Equal(t, []int{1, 2}, res6)
	require.Equal(t, []int{}, res7)
	require.Equal(t, []int{2}, res8)
}

func TestSliceExcludeAll(t *testing.T) {
	t.Parallel()

	res1 := SliceExcludeAll([]int{1, 2, 3})
	res2 := SliceExcludeAll([]int{1, 2, 3}, 0)
	res3 := SliceExcludeAll([]int{1, 2, 3}, 0, 1, 2)
	res4 := SliceExcludeAll([]int{}, 0, 1, 2)

	require.Equal(t, []int{1, 2, 3}, res1)
	require.Equal(t, []int{1, 2, 3}, res2)
	require.Equal(t, []int{3}, res3)
	require.Equal(t, []int{}, res4)
}

func TestSliceExcludeEmpty(t *testing.T) {
	t.Parallel()

	res1 := SliceExcludeEmpty([]int{})
	res2 := SliceExcludeEmpty([]int{0, 1})
	res3 := SliceExcludeEmpty([]int{1, 2, 3})
	res4 := SliceExcludeEmpty([]int{1, 2, 0})

	require.Equal(t, []int{}, res1)
	require.Equal(t, []int{1}, res2)
	require.Equal(t, []int{1, 2, 3}, res3)
	require.Equal(t, []int{1, 2}, res4)

	type a struct {
		A string
	}
	aa := &a{"a"}
	ab := &a{"b"}
	res5 := SliceExcludeEmpty([]*a{aa, nil, ab})
	res6 := SliceExcludeEmpty([]a{{"a"}, {}, {"b"}})

	require.Equal(t, []*a{aa, ab}, res5)
	require.Equal(t, []a{{"a"}, {"b"}}, res6)
}

func TestSliceDifference(t *testing.T) {
	t.Parallel()

	res1, res2 := SliceDifference([]int{}, []int{1, 2})
	res3, res4 := SliceDifference([]int{1, 2}, []int{1, 2})
	res5, res6 := SliceDifference([]int{1, 2}, []int{2, 3})
	res7, res8 := SliceDifference([]int{1, 2}, []int{3, 4})

	require.Equal(t, []int{}, res1)
	require.Equal(t, []int{1, 2}, res2)
	require.Equal(t, []int{}, res3)
	require.Equal(t, []int{}, res4)
	require.Equal(t, []int{1}, res5)
	require.Equal(t, []int{3}, res6)
	require.Equal(t, []int{1, 2}, res7)
	require.Equal(t, []int{3, 4}, res8)
}

func TestSliceUnion(t *testing.T) {
	t.Parallel()

	res1 := SliceUnion[int]()
	res2 := SliceUnion([]int{})
	res3 := SliceUnion([]int{}, []int{})
	res4 := SliceUnion([]int{1}, []int{})
	res5 := SliceUnion([]int{1, 2}, []int{2, 3}, []int{1, 5, 4}, []int{6})

	require.Equal(t, []int{}, res1)
	require.Equal(t, []int{}, res2)
	require.Equal(t, []int{}, res3)
	require.Equal(t, []int{1}, res4)
	require.Equal(t, []int{1, 2, 3, 5, 4, 6}, res5)
}

func TestSliceUnionBy(t *testing.T) {
	t.Parallel()

	res1 := SliceUnionBy(func(_, item int) int {
		return item % 3
	}, []int{0, 1, 2, 3, 4, 5})
	res2 := SliceUnionBy(func(_, item int) int {
		return item % 3
	}, []int{0, 1}, []int{2, 3, 4})

	require.Equal(t, []int{0, 1, 2}, res1)
	require.Equal(t, []int{0, 1, 2}, res2)
}

func TestSliceFilter(t *testing.T) {
	t.Parallel()

	res1 := SliceFilter([]int{0, 1, 2, 3, 4, 5, 6}, func(_ int, item int) bool {
		return item%2 == 0
	})

	require.Equal(t, []int{0, 2, 4, 6}, res1)
}

func TestSliceTransformType(t *testing.T) {
	t.Parallel()

	res1 := SliceTransformType([]int{0, 1, 2}, func(_ int, item int) string {
		return strconv.Itoa(item)
	})

	require.Equal(t, []string{"0", "1", "2"}, res1)
}

func TestSliceFilterTransformType(t *testing.T) {
	t.Parallel()

	res1 := SliceFilterTransformType([]int{0, 1, 2, 3, 4}, func(index int, item int) (string, bool) {
		return strconv.Itoa(item), index%2 == 0
	})

	require.Equal(t, []string{"0", "2", "4"}, res1)
}

func TestSliceFlatten(t *testing.T) {
	t.Parallel()

	res1 := SliceFlatten([][]int{{0, 1}, {2, 3}})

	require.Equal(t, []int{0, 1, 2, 3}, res1)
}

func TestSliceFlattenTransformType(t *testing.T) {
	t.Parallel()

	res1 := SliceFlattenTransformType([]string{"0,1", "2,3"}, func(_ int, item string) []string {
		return strings.Split(item, ",")
	})

	require.Equal(t, []string{"0", "1", "2", "3"}, res1)
}

func TestSliceReduce(t *testing.T) {
	t.Parallel()

	res1 := SliceReduce([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
		return agg + item
	}, 0)
	res2 := SliceReduce([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
		return agg + item
	}, 11)

	require.Equal(t, 10, res1)
	require.Equal(t, 21, res2)
}

func TestSliceReduceRight(t *testing.T) {
	t.Parallel()

	res1 := SliceReduceRight([][]int{{1}, {2}, {4, 3}}, func(agg []int, item []int, _ int) []int {
		return append(agg, item...)
	}, []int{})

	require.Equal(t, []int{4, 3, 2, 1}, res1)
}

func TestSliceGroupBy(t *testing.T) {
	t.Parallel()

	res1 := SliceGroupBy([]int{1, 2, 3, 4}, func(item int) int {
		return item % 3
	})

	require.Equal(t, 3, len(res1))
	require.Equal(t, []int{3}, res1[0])
	require.Equal(t, []int{1, 4}, res1[1])
	require.Equal(t, []int{2}, res1[2])
}

func TestSliceOrderedGroupBy(t *testing.T) {
	t.Parallel()

	res1 := SliceOrderedGroupBy([]int{1, 2, 3, 4}, func(item int) int {
		return item % 3
	})

	require.Equal(t, 3, len(res1))
	require.Equal(t, []int{1, 4}, res1[0])
	require.Equal(t, []int{2}, res1[1])
	require.Equal(t, []int{3}, res1[2])
}

func TestSliceCutChunks(t *testing.T) {
	t.Parallel()

	res1 := SliceCutChunks([]int{1, 2, 3, 4}, 2)
	res2 := SliceCutChunks([]int{1, 2, 3, 4, 5}, 2)

	require.Equal(t, [][]int{{1, 2}, {3, 4}}, res1)
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, res2)
}

func TestSliceInterleaveFlatten(t *testing.T) {
	t.Parallel()

	res1 := SliceInterleaveFlatten([][]int{{0, 1}, {2, 3, 4, 5}, {6, 7, 8}}...)

	require.Equal(t, []int{0, 2, 6, 1, 3, 7, 4, 8, 5}, res1)
}

func TestSliceShuffle(t *testing.T) {
	t.Parallel()

	res1 := SliceShuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	res2 := SliceShuffle([]int{})

	require.NotEqual(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, res1)
	require.Equal(t, []int{}, res2)
}

func TestSliceReverse(t *testing.T) {
	t.Parallel()

	res1 := SliceReverse([]int{0, 1, 2, 3, 4, 5})
	res2 := SliceReverse([]int{0, 1, 2, 1, 2})
	res3 := SliceReverse([]int{})

	require.Equal(t, []int{5, 4, 3, 2, 1, 0}, res1)
	require.Equal(t, []int{2, 1, 2, 1, 0}, res2)
	require.Equal(t, []int{}, res3)
}

func TestSliceFill(t *testing.T) {
	t.Parallel()

	res1 := SliceFill([]foo{"a", "a"}, "b")
	res2 := SliceFill([]foo{}, "b")

	require.Equal(t, []foo{"b", "b"}, res1)
	require.Equal(t, []foo{}, res2)
}

func TestSliceRepeat(t *testing.T) {
	t.Parallel()

	res1 := SliceRepeat[foo](2, "b")
	res2 := SliceRepeat[foo](0, "b")

	require.Equal(t, []foo{"b", "b"}, res1)
	require.Equal(t, []foo{}, res2)
}

func TestSliceRepeatBy(t *testing.T) {
	t.Parallel()

	res1 := SliceRepeatBy(2, func(index int) foo {
		return foo(strconv.Itoa(index))
	})
	res2 := SliceRepeatBy(0, func(index int) foo {
		return foo(strconv.Itoa(index))
	})

	require.Equal(t, []foo{"0", "1"}, res1)
	require.Equal(t, []foo{}, res2)
}

func TestSliceToMap(t *testing.T) {
	t.Parallel()

	res1 := SliceToMap([]string{"s", "sss", "ss"}, func(item string) (int, string) {
		return len(item), item
	})

	require.Equal(t, map[int]string{1: "s", 2: "ss", 3: "sss"}, res1)
}

func TestSliceCutLeft(t *testing.T) {
	t.Parallel()

	require.Equal(t, []int{1, 2, 3, 4}, SliceCutLeft([]int{0, 1, 2, 3, 4}, 1))
	require.Equal(t, []int{2, 3, 4}, SliceCutLeft([]int{0, 1, 2, 3, 4}, 2))
	require.Equal(t, []int{3, 4}, SliceCutLeft([]int{0, 1, 2, 3, 4}, 3))
	require.Equal(t, []int{4}, SliceCutLeft([]int{0, 1, 2, 3, 4}, 4))
	require.Equal(t, []int{}, SliceCutLeft([]int{0, 1, 2, 3, 4}, 5))
	require.Equal(t, []int{}, SliceCutLeft([]int{0, 1, 2, 3, 4}, 6))
}

func TestSliceCutRight(t *testing.T) {
	t.Parallel()

	require.Equal(t, []int{0, 1, 2, 3}, SliceCutRight([]int{0, 1, 2, 3, 4}, 1))
	require.Equal(t, []int{0, 1, 2}, SliceCutRight([]int{0, 1, 2, 3, 4}, 2))
	require.Equal(t, []int{0, 1}, SliceCutRight([]int{0, 1, 2, 3, 4}, 3))
	require.Equal(t, []int{0}, SliceCutRight([]int{0, 1, 2, 3, 4}, 4))
	require.Equal(t, []int{}, SliceCutRight([]int{0, 1, 2, 3, 4}, 5))
	require.Equal(t, []int{}, SliceCutRight([]int{0, 1, 2, 3, 4}, 6))
}

func TestSliceCutLeftOn(t *testing.T) {
	t.Parallel()

	res1 := SliceCutLeftOn([]int{0, 1, 2, 3}, func(item int) bool {
		return item == 2
	})
	res2 := SliceCutLeftOn([]int{}, func(item int) bool {
		return item == 2
	})

	require.Equal(t, []int{2, 3}, res1)
	require.Equal(t, []int{}, res2)
}

func TestSliceCutRightOn(t *testing.T) {
	t.Parallel()

	res1 := SliceCutRightOn([]int{0, 1, 2, 3}, func(item int) bool {
		return item == 2
	})
	res2 := SliceCutRightOn([]int{}, func(item int) bool {
		return item == 2
	})

	require.Equal(t, []int{0, 1, 2}, res1)
	require.Equal(t, []int{}, res2)
}

func TestSliceValueCount(t *testing.T) {
	t.Parallel()

	res1 := SliceValueCount([]int{0, 1, 2, 1, 3, 4}, 1)
	res2 := SliceValueCount([]int{0, 1, 2, 1, 3, 4}, 5)
	res3 := SliceValueCount([]int{}, 5)

	require.Equal(t, 2, res1)
	require.Equal(t, 0, res2)
	require.Equal(t, 0, res3)
}

func TestSliceValueCountBy(t *testing.T) {
	t.Parallel()

	res1 := SliceValueCountBy([]int{0, 1, 2, 1, 3, 4}, func(item int) bool {
		return item == 1
	})
	res2 := SliceValueCountBy([]int{0, 1, 2, 1, 3, 4}, func(item int) bool {
		return item == 5
	})
	res3 := SliceValueCountBy([]int{}, func(item int) bool {
		return item == 1
	})

	require.Equal(t, 2, res1)
	require.Equal(t, 0, res2)
	require.Equal(t, 0, res3)
}

func TestSliceValuesCount(t *testing.T) {
	t.Parallel()

	res1 := SliceValuesCount([]int{0, 1, 1, 2, 2, 2, 3, 3, 3, 3})
	res2 := SliceValuesCount([]int{})

	require.Equal(t, map[int]int{0: 1, 1: 2, 2: 3, 3: 4}, res1)
	require.Equal(t, map[int]int{}, res2)
}

func TestSliceValuesCountBy(t *testing.T) {
	t.Parallel()

	res1 := SliceValuesCountBy([]int{0, 1, 1, 2, 2, 2, 3, 3, 3, 3}, func(item int) int {
		return item % 2
	})
	res2 := SliceValuesCountBy([]int{}, func(item int) int {
		return item % 2
	})

	require.Equal(t, map[int]int{0: 4, 1: 6}, res1)
	require.Equal(t, map[int]int{}, res2)
}

func TestSliceSubset(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5, 6}
	res1 := SliceSubset(arr, 2, 2)
	res2 := SliceSubset(arr, 0, 2)
	res3 := SliceSubset(arr, -3, 2)
	res4 := SliceSubset(arr, 2, 10)
	res5 := SliceSubset(arr, 10, 2)
	res6 := SliceSubset(arr, -10, 2)

	require.Equal(t, []int{3, 4}, res1)
	require.Equal(t, []int{1, 2}, res2)
	require.Equal(t, []int{4, 5}, res3)
	require.Equal(t, []int{3, 4, 5, 6}, res4)
	require.Equal(t, []int{}, res5)
	require.Equal(t, []int{1, 2}, res6)
}

func TestSliceParagraph(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5, 6}
	res1 := SliceParagraph(arr, 0, 0)
	res2 := SliceParagraph(arr, 1, 2)
	res3 := SliceParagraph(arr, 0, 3)
	res4 := SliceParagraph(arr, 1, 7)
	res5 := SliceParagraph(arr, -8, 2)
	res6 := SliceParagraph(arr, 3, 2)
	res7 := SliceParagraph(arr, -5, -2)

	require.Equal(t, []int{}, res1)
	require.Equal(t, []int{2}, res2)
	require.Equal(t, []int{1, 2, 3}, res3)
	require.Equal(t, []int{2, 3, 4, 5, 6}, res4)
	require.Equal(t, []int{1, 2}, res5)
	require.Equal(t, []int{}, res6)
	require.Equal(t, []int{}, res7)
}

func TestSliceReplace(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	res1 := SliceReplace(arr, 3, 5, 0)
	res2 := SliceReplace(arr, 3, 5, 2)
	res3 := SliceReplace(arr, 3, 5, -1)
	res4 := SliceReplace(arr, 6, 5, 1)

	require.Equal(t, arr, res1)
	require.Equal(t, []int{1, 2, 2, 5, 5, 3, 4, 4, 4, 4}, res2)
	require.Equal(t, []int{1, 2, 2, 5, 5, 5, 4, 4, 4, 4}, res3)
	require.Equal(t, arr, res4)
}

func TestSliceReplaceAll(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	res1 := SliceReplaceAll(arr, 3, 5)
	res2 := SliceReplaceAll(arr, 6, 5)

	require.Equal(t, []int{1, 2, 2, 5, 5, 5, 4, 4, 4, 4}, res1)
	require.Equal(t, arr, res2)
}
