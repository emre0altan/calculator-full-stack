FROM node:22-alpine AS web-build
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.21-alpine AS backend-build
WORKDIR /src/backend
COPY backend/ ./
RUN go test ./... && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /calculator-api ./cmd/calculator-api

FROM scratch
WORKDIR /app
COPY --from=backend-build /calculator-api /app/calculator-api
COPY --from=web-build /src/web/dist /app/web
ENV PORT=8080 WEB_DIR=/app/web
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/calculator-api"]
