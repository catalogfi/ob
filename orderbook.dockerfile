FROM golang:1.21-alpine AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go build -tags netgo -ldflags '-s -w' -o ./orderbook ./main.go

RUN chmod +x /app/orderbook

EXPOSE 8080

CMD [ "/app/orderbook" ]


