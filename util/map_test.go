package util

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/rambollwong/rainbowcat/types"
	"github.com/stretchr/testify/require"
)

func TestMapKeys(t *testing.T) {
	t.Parallel()

	res1 := MapKeys(map[string]int{"a": 1, "b": 2})
	sort.Strings(res1)

	require.Equal(t, []string{"a", "b"}, res1)
}

func TestMapValues(t *testing.T) {
	t.Parallel()

	res1 := MapValues(map[string]int{"a": 1, "b": 2})
	sort.Ints(res1)

	require.Equal(t, []int{1, 2}, res1)
}

func TestMapValueOr(t *testing.T) {
	t.Parallel()

	res1 := MapValueOr(map[string]int{"a": 1}, "b", 2)
	require.Equal(t, 2, res1)

	res2 := MapValueOr(map[string]int{"a": 1}, "a", 2)
	require.Equal(t, 1, res2)
}

func TestMapFilter(t *testing.T) {
	t.Parallel()

	res1 := MapFilter(map[string]int{"a": 1, "b": 2, "c": 3}, func(key string, value int) bool {
		return value%2 == 0
	})

	require.Equal(t, map[string]int{"b": 2}, res1)
}

func TestMapFilterByKeys(t *testing.T) {
	t.Parallel()

	res1 := MapFilterByKeys(map[string]int{"a": 1, "b": 2, "c": 3}, []string{"b", "c"})
	require.Equal(t, map[string]int{"b": 2, "c": 3}, res1)

	res2 := MapFilterByKeys(map[string]int{"a": 1, "b": 2, "c": 3}, []string{"d", "e"})
	require.Equal(t, map[string]int{}, res2)
}

func TestMapFilterByValues(t *testing.T) {
	t.Parallel()

	res1 := MapFilterByValues(map[string]int{"a": 1, "b": 2, "c": 3}, []int{2, 3})
	res2 := MapFilterByValues(map[string]int{"a": 1, "b": 2, "c": 2}, []int{2, 3})
	res3 := MapFilterByValues(map[string]int{"a": 1, "b": 2, "c": 3}, []int{2, 2})
	res4 := MapFilterByValues(map[string]int{"a": 1, "b": 2, "c": 3}, []int{4, 5})

	require.Equal(t, map[string]int{"b": 2, "c": 3}, res1)
	require.Equal(t, map[string]int{"b": 2, "c": 2}, res2)
	require.Equal(t, map[string]int{"b": 2}, res3)
	require.Equal(t, map[string]int{}, res4)
}

func TestMapExcludeByKeys(t *testing.T) {
	t.Parallel()

	res1 := MapExcludeByKeys(map[string]int{"a": 1, "b": 2, "c": 3}, []string{})
	res2 := MapExcludeByKeys(map[string]int{"a": 1, "b": 2, "c": 3}, []string{"a", "c"})
	res3 := MapExcludeByKeys(map[string]int{"a": 1, "b": 2, "c": 3}, []string{"d", "e"})
	res4 := MapExcludeByKeys(map[string]int{"a": 1, "b": 2, "c": 3}, []string{"a", "a"})
	res5 := MapExcludeByKeys(map[string]int{}, []string{"a"})

	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, res1)
	require.Equal(t, map[string]int{"b": 2}, res2)
	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, res3)
	require.Equal(t, map[string]int{"b": 2, "c": 3}, res4)
	require.Equal(t, map[string]int{}, res5)
}

func TestMapExcludeByValues(t *testing.T) {
	t.Parallel()

	res1 := MapExcludeByValues(map[string]int{"a": 1, "b": 2, "c": 3}, []int{})
	res2 := MapExcludeByValues(map[string]int{"a": 1, "b": 2, "c": 3}, []int{2})
	res3 := MapExcludeByValues(map[string]int{"a": 1, "b": 2, "c": 3}, []int{2, 2})
	res4 := MapExcludeByValues(map[string]int{"a": 1, "b": 2, "c": 3}, []int{4, 5})
	res5 := MapExcludeByValues(map[string]int{}, []int{2})

	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, res1)
	require.Equal(t, map[string]int{"a": 1, "c": 3}, res2)
	require.Equal(t, map[string]int{"a": 1, "c": 3}, res3)
	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, res4)
	require.Equal(t, map[string]int{}, res5)
}

func TestMapEntries(t *testing.T) {
	t.Parallel()

	res1 := MapEntries(map[string]int{"a": 1, "b": 2, "c": 3})

	sort.Slice(res1, func(i, j int) bool {
		return res1[i].Value < res1[j].Value
	})
	require.Equal(t, []types.Entry[string, int]{
		{
			Key:   "a",
			Value: 1,
		},
		{
			Key:   "b",
			Value: 2,
		},
		{
			Key:   "c",
			Value: 3,
		},
	}, res1)
}

func TestMapFromEntries(t *testing.T) {
	t.Parallel()

	res1 := MapFromEntries([]types.Entry[string, int]{
		{
			Key:   "a",
			Value: 1,
		},
		{
			Key:   "b",
			Value: 2,
		},
		{
			Key:   "c",
			Value: 3,
		},
	})
	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, res1)
}

func TestMapInvert(t *testing.T) {
	t.Parallel()

	res1 := MapInvert(map[string]int{"a": 1, "b": 2, "c": 3})
	require.Equal(t, map[int]string{1: "a", 2: "b", 3: "c"}, res1)
}

func TestMapAssign(t *testing.T) {
	t.Parallel()

	res1 := MapAssign(map[int]int{1: 1}, map[int]int{1: 1, 2: 2}, map[int]int{3: 3})
	require.Equal(t, map[int]int{1: 1, 2: 2, 3: 3}, res1)
}

func TestMapTransformKeys(t *testing.T) {
	t.Parallel()

	res1 := MapTransformKeys(map[string]int{"1": 1, "2": 2}, func(value int, key string) int {
		s, _ := strconv.Atoi(key)
		return s
	})
	require.Equal(t, map[int]int{1: 1, 2: 2}, res1)
}

func TestMapTransformValues(t *testing.T) {
	t.Parallel()

	res1 := MapTransformValues(map[string]int{"1": 1, "2": 2}, func(value int, key string) string {
		return strconv.Itoa(value)
	})
	require.Equal(t, map[string]string{"1": "1", "2": "2"}, res1)
}

func TestMapTransformKeyValues(t *testing.T) {
	t.Parallel()

	res1 := MapTransformKeyValues(map[int]int{1: 1, 2: 2}, func(key int, value int) (float64, string) {
		return float64(key), strconv.Itoa(value)
	})

	require.Equal(t, map[float64]string{1.0: "1", 2.0: "2"}, res1)
}

func TestMapToSlice(t *testing.T) {
	t.Parallel()

	res1 := MapToSlice(map[int]int{1: 2, 2: 3}, func(key int, value int) string {
		return fmt.Sprintf("%d-%d", key, value)
	})
	require.Equal(t, []string{"1-2", "2-3"}, res1)
}
