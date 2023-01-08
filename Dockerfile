FROM golang:1.19.4-alpine3.17

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
