package qfl

// Generates a sequence of numbers from start to end with one step
func generateSequence(start, end int) []int {
	indices := make([]int, end-start)
	for i := 0; i < end-start; i++ {
		indices[i] = start + i
	}

	return indices
}
