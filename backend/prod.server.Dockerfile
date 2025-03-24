# Step 1: Build the Go app
FROM golang:1-alpine3.20 AS builder
WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download
COPY . /app
RUN go build ./cmd/server/main.go -o kessler-server

# Step 2: Export the build result to a plain Alpine image
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/kessler-ingest /app/
COPY . /app
EXPOSE 4041
CMD ["./kessler-server"]
