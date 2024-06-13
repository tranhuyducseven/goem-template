ARG GO_VERSION=1.22.2

# Builder stage
FROM golang:${GO_VERSION} as builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o dist/main cmd/api/main.go

# Use the scratch image as the base image for a small and secure image
FROM debian:12.5-slim
WORKDIR /app

COPY --from=builder /app/dist/ .

ENV PORT=12345
EXPOSE $PORT
ENV GIN_MODE=release

CMD ["/app/main"]
