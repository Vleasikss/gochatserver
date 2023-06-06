FROM golang:1.19-alpine

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY ./static/index.html ./
COPY ./static/main.js ./
COPY ./static/styles.css ./

# Build the Go app
RUN go build -o ./out/go-sample-app socket/gin_melody/main.go

# This container exposes port 8080 to the outside world
EXPOSE 5002

# Run the binary program produced by `go install`
CMD ["./out/go-sample-app"]