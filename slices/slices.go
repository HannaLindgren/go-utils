package slices

import "fmt"

// Insert new element <value> at <index> in <slice>
// @return new slice (old slice is not changed)
func Insert(slice []string, index int, value string) []string {
	if index > len(slice) { // index out of range (index too high)
		panic(fmt.Sprintf("runtime error: index out of range [%v] with length %v (maximum insert index is [%v])",
			index, len(slice), len(slice)))
	}
	if len(slice) == index { // nil or empty slice or after last element
		return append(slice, value)
	}
	slice = append(slice[:index+1], slice[index:]...) // index < len(slice)
	slice[index] = value
	return slice
}
