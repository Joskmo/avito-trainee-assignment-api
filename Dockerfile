FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Install goose for migrations
RUN apk add --no-cache curl
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

EXPOSE 8080

CMD ["./main"]
