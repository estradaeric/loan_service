FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git


ENV GOPROXY=https://proxy.golang.org,direct

# Copy full project
COPY . .
COPY .env.dev .env.dev
COPY .env.prod .env.prod

# Run go mod tidy to generate go.sum in container
RUN go mod tidy

# Build the binary
RUN go build -o loan-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/
COPY --from=builder /app/loan-service .
COPY --from=builder /app/.env.dev .env.dev
COPY --from=builder /app/.env.prod .env.prod

CMD ["./loan-service"]