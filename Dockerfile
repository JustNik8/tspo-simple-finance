FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main ./cmd/web/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["/app/main"]
