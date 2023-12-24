# Stage 1:
FROM golang:1.21.4 AS builder

WORKDIR /app

COPY . .
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o app ./transactions

# Stage 2:
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/transactions/.env .

EXPOSE 8083
CMD ["./app"]