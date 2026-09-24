package httpapi

import (
	"bytes"
	_ "embed"
	"fmt"

	"calculator/backend/internal/calculator"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed calculation.schema.json
var calculationSchemaJSON []byte

var calculationSchema = mustCompileCalculationSchema(calculationSchemaJSON, "calculation.schema.json")

func mustCompileCalculationSchema(source []byte, resource string) *jsonschema.Schema {
	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(source))
	if err != nil {
		panic(fmt.Errorf("parse calculation schema: %w", err))
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(resource, document); err != nil {
		panic(fmt.Errorf("add calculation schema: %w", err))
	}
	schema, err := compiler.Compile(resource)
	if err != nil {
		panic(fmt.Errorf("compile calculation schema: %w", err))
	}
	return schema
}

func validateCalculationRequest(document any, input calculationRequest) string {
	if err := calculationSchema.Validate(document); err == nil {
		return ""
	}

	// Preserve the endpoint's existing error messages for its established rules.
	if input.Operation == "" {
		return "missing required operation"
	}
	op := calculator.Operation(input.Operation)
	if !op.Valid() {
		return "unknown operation"
	}
	if input.A == nil {
		return "missing required number: a"
	}
	if op.NeedsB() && input.B == nil {
		return "missing required number: b"
	}
	return "invalid JSON request body"
}
