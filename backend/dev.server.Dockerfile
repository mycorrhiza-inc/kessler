# run as a nonppriveledged user
FROM golang:1.24 as go-mods
RUN useradd -ms /bin/sh -u 1001 app
USER app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

from go-mods as go-builds

COPY ./internal ./internal
COPY ./cmd/server/ ./cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" ./cmd/

