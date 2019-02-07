FROM golang:1.8 as builder
WORKDIR /usr/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o server .

EXPOSE 80
CMD ["./server"]

