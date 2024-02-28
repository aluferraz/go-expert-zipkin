FROM golang:1.21
WORKDIR /app
RUN go mod tidy
