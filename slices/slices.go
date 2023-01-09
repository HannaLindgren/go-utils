package slices

// Insert new element <value> at <index> in <slice>
func Insert(slice []string, index int, value string) []string {
	if len(slice) == index { // nil or empty slice or after last element
		return append(slice, value)
	}
	slice = append(slice[:index+1], slice[index:]...) // index < len(slice)
	slice[index] = value
	return slice
}
