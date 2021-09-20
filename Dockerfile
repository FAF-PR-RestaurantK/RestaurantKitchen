FROM golang:1.16-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Volcov Oleg <bojikch15@gmail.com>"

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main ./src/main.go

EXPOSE 56566

CMD ["./main"]