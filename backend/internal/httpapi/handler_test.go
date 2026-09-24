package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFallbackHandler(t *testing.T) {
	handler := NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/index.html" {
			t.Fatalf("fallback path = %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/index.html", nil))
	if response.Code != http.StatusAccepted {
		t.Fatalf("fallback status = %d", response.Code)
	}
}

func TestHealthRejectsWrongMethod(t *testing.T) {
	response := httptest.NewRecorder()
	NewHandler(nil).ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/health", nil))
	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("status = %d, Allow = %q", response.Code, response.Header().Get("Allow"))
	}
}

func TestCalculateRejectsUnsupportedMediaTypeAndMalformedJSON(t *testing.T) {
	for _, test := range []struct {
		name, body, contentType string
		status                  int
	}{
		{"unsupported media type", `{}`, "text/plain", http.StatusUnsupportedMediaType},
		{"malformed JSON", `{`, "application/json", http.StatusBadRequest},
		{"oversized body", strings.Repeat(" ", 1<<20) + `{}`, "application/json", http.StatusBadRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/calculate", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()
			NewHandler(nil).ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d, body = %s", response.Code, test.status, response.Body.String())
			}
		})
	}
}

func TestCalculateHandlesDocumentParserFailure(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/calculate", strings.NewReader(`{"operation":"add","a":1,"b":2}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	called := false
	handleCalculate(response, request, func(io.Reader) (any, error) {
		called = true
		return nil, errors.New("parse failed")
	})
	if !called || response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "invalid JSON request body") {
		t.Fatalf("parser called = %t, status = %d, body = %q", called, response.Code, response.Body.String())
	}
}

func TestWriteJSONLogsEncodingFailure(t *testing.T) {
	var logs bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(originalOutput) })
	response := httptest.NewRecorder()
	writeJSON(response, http.StatusOK, make(chan int))
	if response.Code != http.StatusOK || !strings.Contains(logs.String(), "write response: json: unsupported type") {
		t.Fatalf("status = %d, log = %q", response.Code, logs.String())
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name string
		body string
		want float64
	}{
		{"add", `{"operation":"add","a":10,"b":2}`, 12},
		{"subtract", `{"operation":"subtract","a":10,"b":2}`, 8},
		{"multiply", `{"operation":"multiply","a":10,"b":2}`, 20},
		{"divide", `{"operation":"divide","a":10,"b":2}`, 5},
		{"modulus", `{"operation":"modulus","a":10,"b":3}`, 1},
		{"power", `{"operation":"power","a":10,"b":2}`, 100},
		{"sqrt", `{"operation":"sqrt","a":9}`, 3},
		{"sqrt with optional b", `{"operation":"sqrt","a":9,"b":2}`, 3},
		{"percentage", `{"operation":"percentage","a":10,"b":200}`, 20},
		{"fractional negative modulus", `{"operation":"modulus","a":-5.5,"b":2}`, -1.5},
	}

	handler := NewHandler(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/calculate", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			var got calculationResponse
			if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if math.Abs(got.Result-tt.want) > 1e-9 {
				t.Fatalf("result = %v, want %v", got.Result, tt.want)
			}
		})
	}
}

func TestInvalidCalculateRequests(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		method    string
		status    int
		wantError string
	}{
		{"addition overflow", `{"operation":"add","a":1e308,"b":1e308}`, http.MethodPost, http.StatusBadRequest, "addition overflow"},
		{"multiplication overflow", `{"operation":"multiply","a":1e308,"b":2}`, http.MethodPost, http.StatusBadRequest, "multiplication overflow"},
		{"division by zero", `{"operation":"divide","a":10,"b":0}`, http.MethodPost, http.StatusBadRequest, "division by zero"},
		{"division by negative zero", `{"operation":"divide","a":10,"b":-0}`, http.MethodPost, http.StatusBadRequest, "division by zero"},
		{"division overflow", `{"operation":"divide","a":1e308,"b":1e-308}`, http.MethodPost, http.StatusBadRequest, "division overflow"},
		{"modulus by zero", `{"operation":"modulus","a":10,"b":0}`, http.MethodPost, http.StatusBadRequest, "modulus by zero"},
		{"modulus by negative zero", `{"operation":"modulus","a":10,"b":-0}`, http.MethodPost, http.StatusBadRequest, "modulus by zero"},
		{"modulus missing divisor", `{"operation":"modulus","a":10}`, http.MethodPost, http.StatusBadRequest, "missing required number: b"},
		{"modulus missing dividend", `{"operation":"modulus","b":3}`, http.MethodPost, http.StatusBadRequest, "missing required number: a"},
		{"modulus null divisor", `{"operation":"modulus","a":10,"b":null}`, http.MethodPost, http.StatusBadRequest, "missing required number: b"},
		{"modulus invalid operand", `{"operation":"modulus","a":"10","b":3}`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"modulus out of range operand", `{"operation":"modulus","a":1e309,"b":3}`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"invalid numeric input", `{"operation":"add","a":1e309,"b":2}`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"negative square root", `{"operation":"sqrt","a":-1}`, http.MethodPost, http.StatusBadRequest, ""},
		{"non-finite power", `{"operation":"power","a":-1,"b":0.5}`, http.MethodPost, http.StatusBadRequest, ""},
		{"missing operand", `{"operation":"add","a":10}`, http.MethodPost, http.StatusBadRequest, "missing required number: b"},
		{"missing operation", `{"a":10,"b":2}`, http.MethodPost, http.StatusBadRequest, "missing required operation"},
		{"empty object", `{}`, http.MethodPost, http.StatusBadRequest, "missing required operation"},
		{"null operation", `{"operation":null,"a":10,"b":2}`, http.MethodPost, http.StatusBadRequest, "missing required operation"},
		{"unknown operation", `{"operation":"cube","a":10}`, http.MethodPost, http.StatusBadRequest, "unknown operation"},
		{"unknown field", `{"operation":"add","a":10,"b":2,"c":3}`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"wrong top-level type", `[]`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"wrong operation type", `{"operation":true,"a":10,"b":2}`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"wrong second operand type", `{"operation":"add","a":10,"b":"2"}`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"null optional second operand", `{"operation":"sqrt","a":9,"b":null}`, http.MethodPost, http.StatusBadRequest, "invalid JSON request body"},
		{"trailing JSON", `{"operation":"add","a":10,"b":2} {}`, http.MethodPost, http.StatusBadRequest, "request body must contain one JSON object"},
		{"wrong method", ``, http.MethodGet, http.StatusMethodNotAllowed, "method not allowed"},
	}

	handler := NewHandler(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/calculate", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d, body = %s", response.Code, tt.status, response.Body.String())
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("unexpected Content-Type: %s", response.Header().Get("Content-Type"))
			}
			var got errorResponse
			if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil || got.Error == "" {
				t.Fatalf("invalid error response: %s", response.Body.String())
			}
			if tt.wantError != "" && got.Error != tt.wantError {
				t.Fatalf("error = %q, want %q", got.Error, tt.wantError)
			}
		})
	}
}

func TestOldOperationRouteIsGone(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/add", strings.NewReader(`{"a":10,"b":2}`))
	response := httptest.NewRecorder()
	NewHandler(nil).ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
}

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	NewHandler(nil).ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected health response: %d %s", response.Code, response.Body.String())
	}
}
