FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main ./cmd/web/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
RUN chmod 777 /app/main
RUN ls -a /app

EXPOSE 8000
CMD ["/app/main"]