FROM golang:1.22.1-bookworm

ENV APP_HOME /nmbot
WORKDIR "$APP_HOME"
COPY . .
# COPY go.mod go.sum .
RUN go mod download
RUN go build -o main .
CMD ["./main"]
