FROM golang:1.24.3-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .
RUN go build -o wishlist-bot /build/cmd/wishlist-bot
EXPOSE 8080
CMD ["/build/wishlist-bot"]

