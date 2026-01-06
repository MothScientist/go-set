package set

func getCountWorkers(lenSet int) int {
	workers := min(countWorkers, lenSet)
	return workers
}

func getZeroValue[T comparable]() T {
    var zero T
    return zero
}