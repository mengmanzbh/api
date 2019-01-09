FROM golang:1.8 as builder

RUN go get -u github.com/gin-gonic/gin
RUN go get github.com/streadway/amqp
RUN go get github.com/alecthomas/template
RUN go get github.com/swaggo/swag

WORKDIR /go/src/


COPY . /go/src/

RUN go build -o app trainTickets/main.go

# # Application image.
# FROM golang:1.8

# COPY --from=builder /go/src/app/app /usr/local/bin/app

# CMD ["/usr/local/bin/app"]
