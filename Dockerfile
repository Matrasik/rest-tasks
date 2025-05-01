FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rest-app ./cmd/TodoApi/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rest-app .
COPY configs/ configs/
COPY internal/db/migrations internal/db/migrations

CMD ["/app/rest-app"]