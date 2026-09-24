package calculator

import (
	"errors"
	"math"
	"sort"
)

// Operation identifies a supported calculator operation.
type Operation string

const (
	Add        Operation = "add"
	Subtract   Operation = "subtract"
	Multiply   Operation = "multiply"
	Divide     Operation = "divide"
	Modulus    Operation = "modulus"
	Power      Operation = "power"
	SquareRoot Operation = "sqrt"
	Percentage Operation = "percentage"
)

type operationDefinition struct {
	apply  func(float64, float64) (float64, error)
	needsB bool
}

var operations = map[Operation]operationDefinition{
	Add:        {add, true},
	Subtract:   {subtract, true},
	Multiply:   {multiply, true},
	Divide:     {divide, true},
	Modulus:    {modulus, true},
	Power:      {power, true},
	SquareRoot: {squareRoot, false},
	Percentage: {percentage, true},
}

// Operations returns a sorted copy of the supported operation identifiers.
func Operations() []Operation {
	result := make([]Operation, 0, len(operations))
	for op := range operations {
		result = append(result, op)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

// Valid reports whether the operation is supported.
func (op Operation) Valid() bool {
	_, ok := operations[op]
	return ok
}

// NeedsB reports whether an operation requires a second operand.
func (op Operation) NeedsB() bool {
	return operations[op].needsB
}

func add(a, b float64) (float64, error) {
	result := a + b
	if math.IsInf(result, 0) {
		return 0, errors.New("addition overflow")
	}
	return result, nil
}

func subtract(a, b float64) (float64, error) {
	return a - b, nil
}

func multiply(a, b float64) (float64, error) {
	result := a * b
	if math.IsInf(result, 0) {
		return 0, errors.New("multiplication overflow")
	}
	return result, nil
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	result := a / b
	if math.IsInf(result, 0) {
		return 0, errors.New("division overflow")
	}
	return result, nil
}

func modulus(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("modulus by zero")
	}
	return math.Mod(a, b), nil
}

func power(a, b float64) (float64, error) {
	return math.Pow(a, b), nil
}

func squareRoot(a, _ float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("square root requires a non-negative number")
	}
	return math.Sqrt(a), nil
}

func percentage(a, b float64) (float64, error) {
	return a * b / 100, nil
}
