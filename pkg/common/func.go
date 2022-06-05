package common

// Batched returns Two-dimensional slice from the given slice.
func Batched(slice []string, batchSize int) [][]string {
	ret := make([][]string, 0)

	for i := 0; i < len(slice); i += batchSize {
		j := i + batchSize
		if j > len(slice) {
			j = len(slice)
		}

		ret = append(ret, slice[i:j])
	}

	return ret
}
