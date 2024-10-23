FROM golang:1-alpine3.20 
WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod tidy
COPY . /app
EXPOSE 4041
