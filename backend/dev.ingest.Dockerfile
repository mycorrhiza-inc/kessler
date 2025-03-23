# dev dockerfile for ingest pipeline
FROM golang:1.24 AS go-mods
# run as a nonpriveledged user
RUN useradd -ms /bin/sh -u 1001 app
USER app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM go-mods AS go-ingest-builds
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" ./cmd/ingest
