package util

import (
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
