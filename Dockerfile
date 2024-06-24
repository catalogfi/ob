FROM golang:1.21-alpine as builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -tags netgo -ldflags '-s -w' -o ./orderbook ./main.go

FROM alpine:latest  
RUN mkdir store
COPY --from=builder store/setup.sql store/setup.sql
COPY --from=builder local-config.json config.json
COPY --from=builder /app/orderbook    .
EXPOSE 8080
CMD ["./orderbook"]