FROM golang:alpine
RUN apk update && apk --no-cache add autoconf curl git
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
EXPOSE 4004

CMD air -c .air.toml

# FROM golang:alpine

# # Set working directory
# WORKDIR /app

# # Download and install go air live reload (development purpose)
# COPY . ./

# # Command which applies when container from this image runs
# RUN go mod tidy

# # -o = output directory
# RUN go build -o startup main.go

# EXPOSE 4004

# CMD ["./startup"]