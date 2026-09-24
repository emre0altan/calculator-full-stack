package calculator

import (
	"errors"
	"math"
)

// Calculate applies op to a and b. Percentage means a percent of b.
func Calculate(op Operation, a, b float64) (float64, error) {
	if math.IsNaN(a) || math.IsInf(a, 0) {
		return 0, errors.New("a must be a finite number")
	}
	if math.IsNaN(b) || math.IsInf(b, 0) {
		return 0, errors.New("b must be a finite number")
	}

	definition, ok := operations[op]
	if !ok {
		return 0, errors.New("unknown operation")
	}

	result, err := definition.apply(a, b)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, errors.New("result is not a finite number")
	}
	return result, nil
}
