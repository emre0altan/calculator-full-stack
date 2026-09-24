package httpapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"calculator/backend/internal/calculator"
)

func TestSchemaCompilationErrors(t *testing.T) {
	for _, test := range []struct {
		name, source, resource, want string
	}{
		{"invalid JSON", `{`, "calculation.schema.json", "parse calculation schema:"},
		{"invalid resource URL", `{}`, "http://[invalid", "add calculation schema:"},
		{"invalid schema", `{"type":"not-a-type"}`, "calculation.schema.json", "compile calculation schema:"},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				got := recover()
				if got == nil || !strings.Contains(fmt.Sprint(got), test.want) {
					t.Fatalf("panic = %v, want %q", got, test.want)
				}
			}()
			mustCompileCalculationSchema([]byte(test.source), test.resource)
		})
	}
}

func TestSchemaMatchesCalculatorOperations(t *testing.T) {
	var schema struct {
		Properties struct {
			Operation struct {
				Enum []calculator.Operation `json:"enum"`
			} `json:"operation"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(calculationSchemaJSON, &schema); err != nil {
		t.Fatal(err)
	}
	allowed := schema.Properties.Operation.Enum
	sort.Slice(allowed, func(i, j int) bool { return allowed[i] < allowed[j] })
	if !reflect.DeepEqual(allowed, calculator.Operations()) {
		t.Fatalf("schema operations %v differ from calculator operations %v", allowed, calculator.Operations())
	}

	for _, op := range calculator.Operations() {
		t.Run(string(op), func(t *testing.T) {
			document := map[string]any{"operation": string(op), "a": float64(4)}
			validWithoutB := calculationSchema.Validate(document) == nil
			if validWithoutB == op.NeedsB() {
				t.Fatalf("schema and calculator disagree about second operand")
			}
			document["b"] = float64(2)
			if err := calculationSchema.Validate(document); err != nil {
				t.Fatalf("schema rejects valid operands: %v", err)
			}
			if _, err := calculator.Calculate(op, 4, 2); err != nil {
				t.Fatalf("registered operation cannot calculate: %v", err)
			}
		})
	}
}
