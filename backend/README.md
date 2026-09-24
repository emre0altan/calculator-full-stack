# Calculator API

A small REST API written with Go's standard library. It listens on port `8080` by default; set `PORT` to use another port.
Set `WEB_DIR` to a built web directory to serve the UI and API from the same port, as the root Dockerfile does.

## Run locally

Install Go 1.21 or newer. From this `backend` directory, run:

```sh
go run ./cmd/calculator-api
```

The API is available at `http://localhost:8080`. To use a different port, set `PORT` before starting the server (for example, `PORT=9090 go run ./cmd/calculator-api` in a POSIX shell or `$env:PORT=9090; go run ./cmd/calculator-api` in PowerShell). The web development proxy expects port 8080 unless you update [`web/vite.config.ts`](../web/vite.config.ts).

To serve a built web UI from the same server, set `WEB_DIR` to the path of the web app's `dist` directory. For the ready-to-run combined setup, see the [root Docker instructions](../README.md#docker-web-ui-and-api-together).

## Tests and coverage

From this directory, run all backend tests and create a coverage profile with:

```sh
go test ./...
go test ./... '-coverprofile=coverage.out'
go tool cover '-func=coverage.out'
go tool cover '-html=coverage.out' -o coverage.html
```

The last command writes an HTML report to `backend/coverage.html`. Open it in a browser to inspect coverage by source line. The terminal summary and HTML report use the same `coverage.out` profile.

## Project layout

| Folder | Purpose |
| --- | --- |
| `cmd/calculator-api` | Starts the HTTP server |
| `internal/calculator` | Operation types, calculations, and math tests |
| `internal/httpapi` | Routes, JSON Schema request validation, response types, and endpoint tests |

`internal/httpapi.NewHandler` owns API route registration and accepts an optional fallback handler for static files. The server entry point supplies that fallback when `WEB_DIR` is set; new API routes only need to be registered in `internal/httpapi`.

The calculator operation registry supplies execution functions and operand requirements. When adding an operation, update the registry, request schema, and `web/src/operations.ts`. Backend and frontend contract tests catch mismatches. Use `gofmt` for Go formatting and run `go vet ./...` alongside `go test ./...`.

## Endpoints

Send a JSON `POST` request to `/calculate` with `Content-Type: application/json`. Put the operation and operands in the request body:

| Operation | Request body | Calculation |
| --- | --- | --- |
| `add` | `{"operation": "add", "a": 10, "b": 2}` | `a + b` |
| `subtract` | `{"operation": "subtract", "a": 10, "b": 2}` | `a - b` |
| `multiply` | `{"operation": "multiply", "a": 10, "b": 2}` | `a × b` |
| `divide` | `{"operation": "divide", "a": 10, "b": 2}` | `a ÷ b` |
| `modulus` | `{"operation": "modulus", "a": 10, "b": 3}` | Floating-point remainder of `a ÷ b` (result: `1`) |
| `power` | `{"operation": "power", "a": 10, "b": 2}` | `a` raised to `b` |
| `sqrt` | `{"operation": "sqrt", "a": 9}` | Square root of `a` |
| `percentage` | `{"operation": "percentage", "a": 10, "b": 200}` | `a` percent of `b` (result: `20`) |

Successful responses use the form `{"result": 12}`. Invalid requests return `{"error": "..."}` with an appropriate HTTP error status. The backend validates `/calculate` bodies against [`calculation.schema.json`](internal/httpapi/calculation.schema.json): the body must be one JSON object with a supported operation and numeric `a`; all operations except `sqrt` require numeric `b`. Unknown properties and `null` operands are rejected. The schema is compiled once at startup using a Go JSON Schema validator, similar to AJV in a JavaScript service.

Inputs must be finite numbers. Addition, multiplication, and division overflow return specific errors; division or modulus by zero, negative square roots, and other non-finite results are also rejected. Modulus uses Go's `math.Mod`, so the remainder has the sign of `a` and decimal operands are allowed.

`GET /health` returns `{"status":"ok"}`.

Example:

```sh
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation": "add", "a": 10, "b": 2}'
```
