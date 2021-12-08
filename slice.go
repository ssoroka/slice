package slice

import (
	"fmt"
	"sort"
)

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

// Unique returns a new slice that is sorted with all the duplicate strings removed.
func Unique[T Ordered](ss []T) []T {
	return SortedUnique[T](Sort(ss))
}

func SortedUnique[T Ordered](ss []T) []T {
	if ss == nil {
		return nil
	}
	result := []T{}
	last := *new(T)
	for i, s := range ss {
		if i != 0 && last == s {
			continue
		}
		result = append(result, s)
		last = s
	}
	return result
}

// Sort returns a new slice that is the sorted copy of the slice it was called on. Unlike sort.Strings, it does not mutate the original slice
func Sort[T Ordered](ss []T) []T {
	if ss == nil {
		return nil
	}
	ss2 := make([]T, len(ss))
	copy(ss2, ss)
	sort.Slice(ss2, func(i int, j int) bool {
		return ss2[i] <= ss2[j]
	})
	return ss2
}

// SortBy returns a new, slice that is the sorted copy of the slice it was called on, using sortFunc to interpret the string as a sortable integer value. It does not mutate the original slice
func SortBy[T Ordered](ss []T, sortFunc func(slice []T, i, j int) bool) []T {
	if ss == nil {
		return nil
	}
	ss2 := make([]T, len(ss))
	copy(ss2, ss)
	sort.Slice(ss2, func(i, j int) bool {
		return sortFunc(ss2, i, j)
	})
	return ss2
}

// Compare sorts and iterates s1 and s2. calling left() if the element is only in s1, right() if the element is only in s2, and equal() if it's in both.
// this is used as the speedy basis for other set operations.
func Compare[T Ordered](s1, s2 []T, left, equal, right func(s T)) {
	var compareNoop = func(s T) {}
	if left == nil {
		left = compareNoop
	}
	if right == nil {
		right = compareNoop
	}
	if equal == nil {
		equal = compareNoop
	}
	s1 = Unique[T](Sort[T](s1))
	s2 = Unique[T](Sort[T](s2))
	s1Counter := 0
	s2Counter := 0
	for s1Counter < len(s1) && s2Counter < len(s2) {
		if s1[s1Counter] < s2[s2Counter] {
			left(s1[s1Counter])
			s1Counter++
			continue
		}
		if s1[s1Counter] > s2[s2Counter] {
			right(s2[s2Counter])
			s2Counter++
			continue
		}
		// must be equal
		equal(s1[s1Counter])
		s1Counter++
		s2Counter++
	}
	// catch any remaining items
	for i := s1Counter; i < len(s1); i++ {
		left(s1[i])
	}
	for i := s2Counter; i < len(s2); i++ {
		right(s2[i])
	}
}

// Subtract is a set operation that returns the elements from s1 that are not in s2.
func Subtract[T Ordered](s1, s2 []T) []T {
	result := []T{}
	Compare[T](s1, s2, func(s T) {
		result = append(result, s)
	}, nil, nil)
	return result
}

type MapFunc[T any, R any] interface {
	~func(int, T) R | ~func(T) R
}

// Map over each element in the slice and perform an operation on it. the result of the operation will replace the element value.
// Normal func structure is func(i int, s string) string.
// Also accepts func structure func(s string) string
func Map[T any, R any, F MapFunc[T, R]](ss []T, funcInterface F) []R {
	if ss == nil {
		return nil
	}
	f := func(i int, s T) R {
		switch tf := (interface{})(funcInterface).(type) {
		case func(int, T) R:
			return tf(i, s)
		case func(T) R:
			return tf(s)
		}
		panic(fmt.Sprintf("Map cannot understand function type %T", funcInterface))
	}
	result := make([]R, len(ss))
	for i, s := range ss {
		result[i] = f(i, s)
	}
	return result
}

type AccumulatorFunc[T any] func(acc T, i int, s T) T

// Reduce (aka inject) iterates over the slice of items and calls the accumulator function for each pass, storing the state in the acc variable through each pass.
func Reduce[T any](items []T, initialAccumulator T, f AccumulatorFunc[T]) T {
	if items == nil {
		return initialAccumulator
	}
	acc := initialAccumulator
	for i, s := range items {
		acc = f(acc, i, s)
	}
	return acc
}

// Index returns the index of string in the slice, otherwise -1 if the string is not found.
func Index[T comparable](ss []T, s T) int {
	for i, b := range ss {
		if b == s {
			return i
		}
	}
	return -1
}

// SortedIndex returns the index of string in the slice, otherwise -1 if the string is not found.
// this function will do a log2(n) binary search through the list, which is much faster for large lists.
// The slice must be sorted in ascending order.
func SortedIndex[T Ordered](ss []T, s T) int {
	idx := sort.Search(len(ss), func(i int) bool {
		return ss[i] >= s
	})
	if idx >= 0 && idx < len(ss) && ss[idx] == s {
		return idx
	}
	return -1
}

// First returns the First element, or "" if there are no elements in the slice.
// First will also return an "ok" bool value that will be false if there were no elements to select from
func First[T any](ss []T) (T, bool) {
	if len(ss) > 0 {
		return ss[0], true
	}
	return *new(T), false
}

// Last returns the Last element, or "" if there are no elements in the slice.
// Last will also return an "ok" bool value that will be false if there were no elements to select from
func Last[T any](ss []T) (T, bool) {
	if len(ss) > 0 {
		return ss[len(ss)-1], true
	}
	return *new(T), false
}

type SelectFunc[T any] interface {
	~func(int, T) bool | ~func(T) bool
}

func Select[T any, F SelectFunc[T]](ss []T, funcInterface F) []T {
	f := func(i int, s T) bool {
		switch tf := (interface{})(funcInterface).(type) {
		case func(int, T) bool:
			return tf(i, s)
		case func(T) bool:
			return tf(s)
		default:
			panic(fmt.Sprintf("Filter cannot understand function type %T", funcInterface))
		}
	}

	result := []T{}

	for i, s := range ss {
		if f(i, s) {
			result = append(result, s)
		}
	}
	return result
}

// Contains returns true if the string is in the slice.
func Contains[T comparable](ss []T, s T) bool {
	return Index(ss, s) != -1
}

// SortedContains returns true if the string is in an already sorted slice. it's faster than Contains for large slices
func SortedContains[T Ordered](ss []T, s T) bool {
	return SortedIndex(ss, s) != -1
}

// Pop pops the last element off a slice and returns the popped element and the remaining slice
// (note that the original slice is not modified)
func Pop[T any](ss []T) (T, []T) {
	elem, ok := Last(ss)
	if ok {
		return elem, ss[0 : len(ss)-1]
	}

	return elem, nil
}

// Shift returns the first element and the remaining slice
func Shift[T any](ss []T) (T, []T) {
	if len(ss) == 0 {
		return *new(T), nil
	}

	return ss[0], ss[1:]
}

// Unshift prepends the element in front of the first value
func Unshift[T any](ss []T, elem T) []T {
	return append([]T{elem}, ss...)
}

type FindFunc[T any] interface {
	~func(T) bool | ~func(int, T) bool
}

// Find and return the first element that matches. returns false if none found.
func Find[T any, F FindFunc[T]](ss []T, funcInterface F) (elem T, found bool) {
	f := func(i int, s T) bool {
		switch tf := (interface{})(funcInterface).(type) {
		case func(int, T) bool:
			return tf(i, s)
		case func(T) bool:
			return tf(s)
		default:
			panic(fmt.Sprintf("Find cannot understand function type %T", funcInterface))
		}
	}

	for i, s := range ss {
		if f(i, s) {
			return s, true
		}
	}
	return *new(T), false
}
