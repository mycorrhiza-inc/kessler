# run as a nonpriveledged user
FROM golang:1.24 AS go-mods
RUN useradd -ms /bin/sh -u 1001 app
USER app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

from go-mods AS go-builds
WORKDIR /app
COPY ./internal ./internal
COPY ./cmd/server/ ./cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" ./cmd/
