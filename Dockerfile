FROM golang:1.21-alpine 

RUN mkdir /app
WORKDIR /app

COPY . .

ADD local-config.json config.json

RUN go build -tags netgo -ldflags '-s -w' -o ./orderbook ./main.go

EXPOSE 8080

RUN chmod +x /app/orderbook
CMD [ "/app/orderbook" ]