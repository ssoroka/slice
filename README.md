# The missing slice package

A Go-generics (Go 1.18) based functional library with no side-effects that adds the following functions to a `slice` package:

- `Unique()`: removes duplicate items from a slice
- `SortedUnique()`: removes duplicate items from a sorted slice. Fast!
- `Sort()`: Sorts a slice
- `SortBy()`: Sorts a slice by arbitrary criteria
- `Compare()`: Iterates two slices comparing them together
- `Subtract()`: subtract one slice from another
- `Map()`: map over a slice and transform it
- `Reduce()`: reduce a slice to a single value
- `Index()`: return the index of an element
- `SortedIndex()`: return the index of an element in a sorted list. Fast!
- `First()`: return the first element
- `Last()`: return the last element
- `Select()`: select all elements matching some specified criteria
- `Contains()`: returns true if the slice contains an element
- `SortedContains()`: returns true if the sorted slice contains an element. Fast!
- `Pop()`: pop the last element off the list
- `Shift()`: shift the first element off the list
- `Unshift()`: shift an element into the first place (prepend)
- `Find()`: find the first element in a slice that matches some criteria

Method signatures:

(where T is almost any type)

- slice.Unique([]T) []T 
- slice.SortedUnique([]T) []T 
- slice.Sort([]T) []T 
- slice.SortBy([]T, sortFunc func(slice []T, i, j int) bool) []T 
- slice.Compare(s1, s2 []T, left, equal, right func(elem T)) 
- slice.Subtract(s1, s2 []T) []T 
- slice.Map([]T, func(i int, elem T) T) []T 
- slice.Reduce(items []T, initialAccumulator T, f AccumulatorFunc[T]) T 
- slice.Index([]T, elem T) int 
- slice.SortedIndex([]T, elem T) int 
- slice.First([]T) (T, bool) 
- slice.Last([]T) (T, bool) 
- slice.Select([]T, func(i int, elem T) T) []T 
- slice.Contains([]T, elem T) bool 
- slice.SortedContains([]T, elem T) bool 
- slice.Pop([]T) (T, []T) 
- slice.Shift([]T) (T, []T) 
- slice.Unshift([]T, elem T) []T 
- slice.Find([]T, func(i int, elem T) T) (elem T, found bool) 

Note the AccumulatorFunc signature is `func(acc T, i int, elem T) T`

## Examples

See tests for more examples. Here are a few:

```go
  // sum function implemented with Reduce
  input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := slice.Reduce(input, 0, func(acc int, i int, elem int) string {
		return acc + elem
	})
  // result == 55
```

map from complex structs to string slices

```go
type Group struct {
	Name string
}
groups := []Group{{Name: "users"}, {Name: "admins"}}
names := slice.Map[Group, string](groups, func(group Group) string {
	return group.Name
})
fmt.Println(names) // outputs ["users", "admins"]
```
