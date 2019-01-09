FROM golang:1.8 as builder

RUN go get -u github.com/gin-gonic/gin
RUN go get github.com/streadway/amqp
RUN go get github.com/alecthomas/template
RUN go get github.com/swaggo/swag

COPY . /go/src/


WORKDIR /go/src/trainTickets
RUN go build -o main main.go
CMD ["./main"]
# Application image.
# FROM golang:1.8

# COPY --from=builder /go/src/trainTickets/app /usr/local/bin/app

# CMD ["/usr/local/bin/app"]
