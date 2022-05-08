FROM golang:1.12.0-alpine

ENV GO111MODULE=on
RUN mkdir /app
Add . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]