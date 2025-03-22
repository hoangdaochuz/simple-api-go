FROM golang:1.23.0 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gin_app

FROM alpine:latest

WORKDIR /root
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/gin_app .

COPY .env .env

EXPOSE 8080
CMD ["./gin_app"]
