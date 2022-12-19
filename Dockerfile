FROM golang:1.18

WORKDIR /go/src
COPY . .

CMD ["go", "run", "main.go"]
