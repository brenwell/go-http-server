# Use the official Golang image as a base
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o server .

# Use a lightweight Alpine image for the final stage
FROM alpine:latest
RUN addgroup -S -g 3000 appgroup && adduser -S -u 1000 -G appgroup appuser
WORKDIR /app
COPY --from=builder --chown=appuser:appgroup /app/server .
USER 1000
EXPOSE 8080
CMD ["./server"]