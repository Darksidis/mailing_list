FROM golang:1.19-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
COPY database.go ./

RUN apk add gcc
RUN apk add musl-dev
RUN apk add libc-dev

RUN go build -o main

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY .env ./
COPY mail_list.db ./
COPY templates/ /app/templates/

EXPOSE 8000

CMD [ "/app/main" ]

