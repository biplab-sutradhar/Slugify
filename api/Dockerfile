# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

# Migrations are loaded at runtime via file://migrations
COPY --from=builder /out/server /app/server
COPY --from=builder /src/migrations /app/migrations

USER nonroot:nonroot
EXPOSE 9000
ENTRYPOINT ["/app/server"]