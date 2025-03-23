# Step 1: Build the Go app
FROM golang:1-alpine3.20 AS builder
WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download

FROM cosmtrek/air
WORKDIR /app
COPY --from=builder /app/ /app/
COPY ./ingest.air.toml /app/.air.toml
COPY ./cmd/ingest/main.go /app
EXPOSE 4041
CMD ["air"]
