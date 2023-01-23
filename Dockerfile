FROM golang:alpine
RUN apk update && apk --no-cache add autoconf curl git
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
EXPOSE 3939
CMD air -c .air.toml