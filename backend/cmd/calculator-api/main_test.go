package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	originalServe, originalFatal := listenAndServe, fatalf
	t.Cleanup(func() {
		listenAndServe, fatalf = originalServe, originalFatal
	})

	for _, test := range []struct {
		name, port, webDir, wantAddr string
		serveErr                     error
		wantFatal                    string
	}{
		{"default port", "", "", ":8080", nil, ""},
		{"configured port and closed server", "9090", t.TempDir(), ":9090", http.ErrServerClosed, ""},
		{"listen error", "invalid", "", ":invalid", errors.New("listen failed"), "serve calculator API: listen failed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("PORT", test.port)
			t.Setenv("WEB_DIR", test.webDir)
			called := false
			listenAndServe = func(server *http.Server) error {
				called = true
				if server.Addr != test.wantAddr || server.ReadHeaderTimeout != 5*time.Second {
					t.Fatalf("server address = %q, header timeout = %s", server.Addr, server.ReadHeaderTimeout)
				}
				response := httptest.NewRecorder()
				server.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))
				if response.Code != http.StatusOK {
					t.Fatalf("health status = %d", response.Code)
				}
				return test.serveErr
			}
			var fatalMessage string
			fatalf = func(format string, args ...any) {
				fatalMessage = fmt.Sprintf(format, args...)
			}
			main()
			if !called || fatalMessage != test.wantFatal {
				t.Fatalf("server called = %t, fatal message = %q, want %q", called, fatalMessage, test.wantFatal)
			}
		})
	}
}

func TestCombinedHandler(t *testing.T) {
	webDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(webDir, "index.html"), []byte("<h1>Calculator</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := newHandler(webDir)

	for _, test := range []struct {
		name   string
		method string
		path   string
		body   string
		want   string
	}{
		{"web", http.MethodGet, "/", "", "<h1>Calculator</h1>"},
		{"health", http.MethodGet, "/health", "", `{"status":"ok"}`},
		{"calculate", http.MethodPost, "/calculate", `{"operation":"add","a":10,"b":2}`, `{"result":12}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			if test.body != "" {
				request.Header.Set("Content-Type", "application/json")
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), test.want) {
				t.Fatalf("status = %d, body = %q", response.Code, response.Body.String())
			}
		})
	}
}
