FROM golang:1.21-alpine AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go build -tags netgo -ldflags '-s -w' -o ./rest ./cmd/rest/rest.go

RUN go build -tags netgo -ldflags '-s -w' -o ./watcher ./cmd/watcher/watcher.go

# move rest to a smaller docker image
FROM alpine:latest AS rest

COPY --from=builder /app/rest /app

COPY --from=builder /app/config.json /app

COPY --from=builder /app/store/setup.sql /app/store

RUN chmod +x /app/rest

EXPOSE 8080:8080

CMD [ "/app/rest" ]

# move watcher to a smaller docker image
FROM alpine:latest AS watcher

COPY --from=builder /app/watcher /app

COPY --from=builder /app/config.json /app

COPY --from=builder /app/store/setup.sql /app/store

RUN chmod +x /app/watcher

CMD [ "/app/watcher" ]

