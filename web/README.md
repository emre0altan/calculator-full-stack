# Calculator web UI

This React and Vite app sends calculations to the Go API. For a one-container setup of both projects, see the [root README](../README.md#docker-web-ui-and-api-together).

## Run locally

Install a current Node.js 22 or 24 release with npm. Start the [backend API](../backend/README.md#run-locally) on its default port, 8080, in another terminal. Then, from this `web` directory, run:

```sh
npm ci
npm run dev
```

Open the URL printed by Vite (usually http://localhost:5173). Its development server forwards `/calculate` requests to `http://localhost:8080`. If the API is running on another port, update the proxy target in `vite.config.ts` before starting Vite.

## Build and test

From this directory:

```sh
npm run build
npm test
npm run lint
npm run format:check
```

The build writes static files to `dist/`. To serve that build locally with Vite, run `npm run preview` and open the URL it prints. Vite's development proxy is configured for `npm run dev`; use the root Docker setup to serve the built UI with the API on one origin.

Run `npm run check` for all checks, or `npm run format` to apply the shared formatting rules. Contract tests read the backend JSON schema, so run tests from a checkout containing both `web/` and `backend/`.

## Coverage

From this directory, run:

```sh
npm run test:coverage
```

The command prints a coverage summary and refreshes the HTML report at `coverage/index.html`. Open that file to inspect individual source lines. The coverage package is included in `npm ci`, and `npm run check` refreshes the report. Coverage includes frontend source files and excludes test and setup files.

## Project layout

| File                 | Purpose                                                      |
| -------------------- | ------------------------------------------------------------ |
| `src/main.tsx`       | React startup and global styles                              |
| `src/App.tsx`        | Page layout and introduction                                 |
| `src/Calculator.tsx` | Calculator form, state, validation, and feedback             |
| `src/operations.ts`  | Operation metadata and the derived TypeScript operation type |
| `src/api.ts`         | HTTP requests and response handling                          |
| `src/styles.css`     | Page and calculator styles                                   |
| `src/*.test.ts(x)`   | UI, API client, and operation contract tests                 |
| `src/test/setup.ts`  | Shared test setup and cleanup                                |

When adding an operation, update `operations.ts`, the backend calculator registry, and the backend request schema. Contract tests check that supported operations and second-operand requirements stay aligned. Keep feature-specific code together as the app grows; the current single-screen app does not need additional folder layers.
