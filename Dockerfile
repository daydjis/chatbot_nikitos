FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o nikitos-bot .

CMD ["./nikitos-bot"]
