package slice_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/ssoroka/slice"
	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	list := []string{"foo", "bar", "zee"}

	require.True(t, slice.Contains(list, "foo"))
	require.False(t, slice.Contains(list, "toodles"))
}

func TestSubtract(t *testing.T) {
	result := slice.Subtract([]string{"A", "B", "C"}, []string{"C"})
	require.Equal(t, []string{"A", "B"}, result)

	result2 := slice.Subtract([]int{1, 2, 3}, []int{2, 3})
	require.Equal(t, []int{1}, result2)
}

func TestUnique(t *testing.T) {
	require.Equal(t, []string{"A", "B", "C"}, slice.Unique([]string{"A", "B", "C", "A", "B", "C", "B", "C", "A"}))
	require.Equal(t, []int64{1, 2, 3}, slice.Unique([]int64{1, 2, 3, 1, 2, 3, 2, 3, 1}))
}

func TestMap(t *testing.T) {
	s := []string{"a", "b", "c"}
	result := slice.Map(s, func(i int, s string) string {
		return strings.ToUpper(s)
	})
	expected := []string{"A", "B", "C"}

	require.Equal(t, expected, result)

	var nilStringSlice []string
	require.Equal(t, []string(nil), slice.Map(nilStringSlice, nil))
	require.Equal(t, []string{}, slice.Map([]string{}, nil))
	require.Equal(t, []string{"a"}, slice.Map([]string{"a"}, nil))

	require.Equal(t, []string{"FISH"}, slice.Map([]string{"fish"}, strings.ToUpper))
	require.Equal(t, []string{"fish"}, slice.Map([]string{" fish "}, strings.TrimSpace))

}

func TestReduce(t *testing.T) {
	s := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	// sum up strings as if they were ints.
	result := slice.Reduce(s, "0", func(acc string, i int, s string) string {
		accumulator, _ := strconv.Atoi(acc)
		current, _ := strconv.Atoi(s)
		s = strconv.Itoa(accumulator + current)
		return s
	})

	require.Equal(t, "55", result)
}

func TestSelect(t *testing.T) {
	result := slice.Select([]string{"car1", "car2", "bus1", "bus2"}, func(i int, s string) bool {
		return strings.HasPrefix(s, "car")
	})

	require.Equal(t, []string{"car1", "car2"}, result)
}

func TestFirst(t *testing.T) {
	first, ok := slice.First([]int{1, 2, 3})
	require.Equal(t, int(1), first)
	require.Equal(t, true, ok)

	first2, ok := slice.First([]string{})
	require.Equal(t, "", first2)
	require.Equal(t, false, ok)
}

func TestLast(t *testing.T) {
	last, ok := slice.Last([]int{1, 2, 3})
	require.Equal(t, int(3), last)
	require.Equal(t, true, ok)

	last, ok = slice.Last([]int{})
	require.Equal(t, int(0), last)
	require.Equal(t, false, ok)
}

func TestPop(t *testing.T) {
	el, newSlice := slice.Pop([]int8{1, 2, 3, 4})
	require.Equal(t, []int8{1, 2, 3}, newSlice)
	require.Equal(t, int8(4), el)
}

func TestUnshift(t *testing.T) {
	newSlice := slice.Unshift([]int{1, 2, 3, 4}, 5)
	require.Equal(t, []int{5, 1, 2, 3, 4}, newSlice)
}

func TestShift(t *testing.T) {
	el, newSlice := slice.Shift([]int{1, 2, 3, 4})
	require.Equal(t, []int{2, 3, 4}, newSlice)
	require.Equal(t, 1, el)
}

func TestSortBy(t *testing.T) {
	// func SortBy[T Ordered](ss []T, sortFunc func(slice []T, i, j int) bool) []T {

	result := slice.SortBy([]string{"High = 3", "Low = 1", "Nominal = 2"}, func(sl []string, i, j int) bool {
		last, _ := slice.Last(strings.Split(sl[i], " "))
		iVal, _ := strconv.Atoi(last)
		last, _ = slice.Last(strings.Split(sl[j], " "))
		jVal, _ := strconv.Atoi(last)
		return iVal < jVal
	})

	expected := []string{"Low = 1", "Nominal = 2", "High = 3"}

	require.Equal(t, expected, result)
}

func TestFind(t *testing.T) {
	result, found := slice.Find([]int{10, 20, 30, 44, 52, 66, 77, 81, 93, 111}, func(elem int) bool {
		return elem%7 == 0
	})

	require.Equal(t, 77, result)
	require.Equal(t, true, found)

	result, found = slice.Find([]int{1, 2, 3}, func(elem int) bool {
		return elem%102 == 0
	})

	require.Equal(t, 0, result)
	require.Equal(t, false, found)
}
