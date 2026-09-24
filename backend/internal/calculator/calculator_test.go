package calculator

import (
	"math"
	"reflect"
	"testing"
)

func TestOperationsMetadata(t *testing.T) {
	want := []Operation{Add, Divide, Modulus, Multiply, Percentage, Power, SquareRoot, Subtract}
	if got := Operations(); !reflect.DeepEqual(got, want) {
		t.Fatalf("operations = %v, want %v", got, want)
	}
	result := Operations()
	result[0] = "changed"
	if got := Operations(); !reflect.DeepEqual(got, want) {
		t.Fatalf("modifying returned operations changed registry: %v", got)
	}
	for _, op := range want {
		if !op.Valid() || op.NeedsB() != (op != SquareRoot) {
			t.Fatalf("incorrect metadata for %q", op)
		}
	}
	if Operation("unknown").Valid() {
		t.Fatal("unknown operation reported as valid")
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name string
		op   Operation
		a, b float64
		want float64
	}{
		{"addition", Add, 10, 2, 12},
		{"addition at finite limit", Add, math.MaxFloat64, 0, math.MaxFloat64},
		{"subtraction", Subtract, 10, 2, 8},
		{"multiplication", Multiply, 10, 2, 20},
		{"multiplication at finite limit", Multiply, math.MaxFloat64, 1, math.MaxFloat64},
		{"division", Divide, 10, 2, 5},
		{"large finite division", Divide, 1e308, 10, 1e307},
		{"zero numerator", Divide, 0, 2, 0},
		{"modulus", Modulus, 10, 3, 1},
		{"fractional modulus", Modulus, 5.5, 2, 1.5},
		{"negative dividend modulus", Modulus, -10, 3, -1},
		{"negative divisor modulus", Modulus, 10, -3, 1},
		{"zero dividend modulus", Modulus, 0, 3, 0},
		{"power", Power, 10, 2, 100},
		{"square root", SquareRoot, 9, 0, 3},
		{"percentage", Percentage, 10, 200, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.op, tt.a, tt.b)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(got-tt.want) > 1e-9*math.Max(1, math.Abs(tt.want)) {
				t.Fatalf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name string
		op   Operation
		a, b float64
		want string
	}{
		{"addition overflow", Add, math.MaxFloat64, math.MaxFloat64, "addition overflow"},
		{"negative addition overflow", Add, -math.MaxFloat64, -math.MaxFloat64, "addition overflow"},
		{"multiplication overflow", Multiply, math.MaxFloat64, 2, "multiplication overflow"},
		{"negative multiplication overflow", Multiply, -math.MaxFloat64, 2, "multiplication overflow"},
		{"division by zero", Divide, 10, 0, "division by zero"},
		{"division by negative zero", Divide, 10, math.Copysign(0, -1), "division by zero"},
		{"division overflow", Divide, 1e308, 1e-308, "division overflow"},
		{"negative division overflow", Divide, -1e308, 1e-308, "division overflow"},
		{"modulus by zero", Modulus, 10, 0, "modulus by zero"},
		{"modulus by negative zero", Modulus, 10, math.Copysign(0, -1), "modulus by zero"},
		{"invalid modulus a", Modulus, math.Inf(1), 3, "a must be a finite number"},
		{"invalid modulus b", Modulus, 10, math.NaN(), "b must be a finite number"},
		{"invalid a", Add, math.Inf(1), 2, "a must be a finite number"},
		{"invalid b", Multiply, 2, math.NaN(), "b must be a finite number"},
		{"negative square root", SquareRoot, -1, 0, "square root requires a non-negative number"},
		{"non-finite power", Power, -1, 0.5, "result is not a finite number"},
		{"unknown operation", Operation("unknown"), 1, 2, "unknown operation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Calculate(tt.op, tt.a, tt.b); err == nil || err.Error() != tt.want {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestModulusWithLargeQuotient(t *testing.T) {
	result, err := Calculate(Modulus, 1e308, 1e-308)
	if err != nil {
		t.Fatal(err)
	}
	if math.IsNaN(result) || math.IsInf(result, 0) || math.Abs(result) >= 1e-308 {
		t.Fatalf("unexpected modulus result: %v", result)
	}
}
