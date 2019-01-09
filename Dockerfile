FROM golang:1.8 as builder

RUN go get -u github.com/gin-gonic/gin
RUN go get github.com/streadway/amqp

WORKDIR /go/src/app/


COPY . .
# ADD ./trainTickets/. .
# # COPY trainTickets/main.go .
RUN /usr/local/go/bin/go build -o app trainTickets/main.go

# # Application image.
# FROM golang:1.8

# COPY --from=builder /go/src/app/app /usr/local/bin/app

# CMD ["/usr/local/bin/app"]
