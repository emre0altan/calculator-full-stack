# Calculator full stack

## Setup

### Docker (web UI and API together)

From the repository root, build and run the image with Docker:

```sh
docker build -t calculator-full-stack .
docker run --rm -p 8080:8080 calculator-full-stack
```

Open http://localhost:8080 for the calculator UI. The API is available on the same port. Stop the container with Ctrl+C.

### Local development (web UI and API separately)

Install Go 1.21 or newer and a current Node.js 22 or 24 release with npm. Start the API from the repository root in one terminal:

```sh
cd backend
go run ./cmd/calculator-api
```

In a second terminal, start the web app from the repository root:

```sh
cd web
npm ci
npm run dev
```

Open the URL printed by Vite (usually http://localhost:5173). Keep the API running on its default port, 8080: the Vite development server forwards `/calculate` requests there. See the [backend](backend/README.md) and [web](web/README.md) READMEs for folder-specific commands and checks.

## Tests and coverage

Run the backend and web test suites from the repository root:

```sh
cd backend
go test ./...
cd ../web
npm ci
npm test
```

To create backend coverage reports, run `go test ./... '-coverprofile=coverage.out'` from `backend/`, then `go tool cover '-func=coverage.out'` for a terminal summary or `go tool cover '-html=coverage.out' -o coverage.html` for an HTML report. For web coverage, run `npm run test:coverage` from `web/`; this refreshes the [frontend coverage report](web/coverage/index.html).

## API examples

These examples use the API at `http://localhost:8080` with either setup above. The `curl` commands use POSIX shell syntax. In PowerShell, you can make a JSON request with `Invoke-RestMethod -Uri http://localhost:8080/calculate -Method Post -ContentType 'application/json' -Body '{"operation":"add","a":10,"b":2}'`.

Check that the API is running:

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```

Add two numbers:

```sh
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":10,"b":2}'
# {"result":12}
```

Calculate a square root (only `a` is needed):

```sh
curl -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":9}'
# {"result":3}
```

Invalid calculations return an error with HTTP 400, for example division by zero:

```sh
curl -i -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":10,"b":0}'
# HTTP/1.1 400 Bad Request
# {"error":"division by zero"}
```

See the [backend API reference](backend/README.md#endpoints) for all supported operations and request rules.

## Design decisions

- The Go API uses the standard library for HTTP routing and JSON handling. A JSON Schema validator checks calculation requests against a declarative schema before the math runs.
- A single `POST /calculate` endpoint selects the operation from JSON. This keeps the request and response shape consistent across the eight operations.
- Calculation and HTTP handling live in separate backend packages so the math can be tested without an HTTP server. Invalid inputs and non-finite results return JSON errors instead of unusable numeric values.
- The React UI sends requests to `/calculate` on its own origin. Vite proxies that route during local development; the Go server serves the built UI and API together in Docker, so neither mode needs a browser-side API URL setting.
- The Dockerfile builds the UI and API in separate stages and copies only the built assets and server binary into the final image.
