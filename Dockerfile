# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./server/main.go

# Run stage
FROM alpine:latest
WORKDIR /root/

RUN apk --no-cache add ca-certificates 
COPY --from=builder /app/main .
COPY --from=builder /app/internal/migrations ./migrations
EXPOSE 8080
CMD ["./main"]