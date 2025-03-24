# run as a nonpriveledged user
FROM golang:1.24 AS go-mods
RUN useradd -ms /bin/sh -u 1001 app
USER app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM go-mods AS go-server-builds
WORKDIR /app
COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" ./cmd/server
