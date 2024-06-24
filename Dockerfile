FROM golang:1.21-alpine 
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -tags netgo -ldflags '-s -w' -o ./orderbook ./main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
ADD local-config.json config.json
COPY --from=builder /app/orderbook    .
EXPOSE 8080
CMD ["./orderbook"]