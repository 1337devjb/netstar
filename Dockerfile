FROM golang:1.18 as builder

RUN mkdir /app
Add . /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY .env ./
COPY *.go ./

RUN go build -o /docker-gs-ping
CMD [ "/docker-gs-ping" ]
