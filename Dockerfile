FROM golang:1.21-alpine as builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -tags netgo -ldflags '-s -w' -o ./orderbook ./main.go

FROM alpine:latest  
RUN mkdir /app
WORKDIR /app
RUN mkdir store
COPY --from=builder /app/store/setup.sql /app/store/setup.sql
COPY --from=builder /app/local-config.json /app/config.json
COPY --from=builder /app/orderbook    /app/.
EXPOSE 8080
CMD ["/app/orderbook"]