FROM golang:1.21
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN go build -o /telegram-bot
CMD ["/telegram-bot"]