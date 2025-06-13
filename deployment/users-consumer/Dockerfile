FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o user-generator ./cmd/users-consumer/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-generator .

CMD [ "./user-generator" ]