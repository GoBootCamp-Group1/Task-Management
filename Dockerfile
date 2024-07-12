FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build cmd/api/main.go

EXPOSE 8082

CMD ["./main", "--config=/app/cmd/api/config.yaml"]


