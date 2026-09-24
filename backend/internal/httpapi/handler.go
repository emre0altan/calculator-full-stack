package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"

	"calculator/backend/internal/calculator"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// NewHandler registers API routes and an optional fallback for static files.
// A nil fallback leaves unmatched routes as HTTP 404 responses.
func NewHandler(fallback http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/calculate", calculateHandler)
	if fallback != nil {
		mux.Handle("/", fallback)
	}
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	handleCalculate(w, r, jsonschema.UnmarshalJSON)
}

func handleCalculate(w http.ResponseWriter, r *http.Request, parseDocument func(io.Reader) (any, error)) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var body json.RawMessage
	if err := decoder.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return
	}
	document, err := parseDocument(bytes.NewReader(body))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")
		return
	}
	var input calculationRequest
	if err := json.Unmarshal(body, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")
		return
	}
	if message := validateCalculationRequest(document, input); message != "" {
		writeError(w, http.StatusBadRequest, message)
		return
	}

	op := calculator.Operation(input.Operation)
	var b float64
	if input.B != nil {
		b = *input.B
	}
	result, err := calculator.Calculate(op, *input.A, b)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, calculationResponse{Result: result})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}
