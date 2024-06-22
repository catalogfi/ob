FROM golang:1.21-alpine 

RUN mkdir /app

WORKDIR /app

COPY . .

ADD local-config.json config.json

EXPOSE 8080

CMD [ "go", "run", "main.go" ]


