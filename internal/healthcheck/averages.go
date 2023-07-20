package healthcheck

import (
	"golang.org/x/exp/constraints"
	"guarde/pkg/utils"
	"sync"
)

func NewAverage[T constraints.Integer](limit uint8) *Average[T] {
	return &Average[T]{limit: limit}
}

type Average[T constraints.Integer] struct {
	container []T
	limit     uint8
	mutex     sync.Mutex
	stale     int
}

// Stale gets the average number, although it is called stale, it's more than likely to be the
// latest numbers and can be used instead of the GetAverage function which calculates.
func (avg *Average[T]) Stale() int {
	return avg.stale
}

func (avg *Average[T]) Values() []T {
	return utils.ReturningMutex(&avg.mutex, func() []T {
		arr := make([]T, len(avg.container))
		copy(arr, avg.container)
		return arr
	})
}

func (avg *Average[T]) average() int {
	average := 0
	for _, value := range avg.container {
		average += int(value)
	}
	if average == 0 {
		return 0
	}
	return average / len(avg.container)
}

func (avg *Average[T]) GetAverage() int {
	return utils.ReturningMutex(&avg.mutex, func() int {
		return avg.average()
	})
}

func (avg *Average[T]) Add(value T) {
	utils.UseMutex(&avg.mutex, func() {
		occupied := len(avg.container)

		freeAmount := 0
		if occupied > int(avg.limit) {
			freeAmount = occupied - int(avg.limit)
		}

		avg.container = append(avg.container[freeAmount:], value)
		avg.stale = avg.average()
	})
}
