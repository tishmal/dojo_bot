# Этап сборки
FROM golang:1.24.3 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /dojo_bot

# Финальный образ
FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /dojo_bot /dojo_bot
CMD ["/dojo_bot"]