FROM golang:1.20-alpine3.17 as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o apiApp ./cmd/api

RUN chmod +x /app/apiApp

#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/apiApp /app

CMD ["/app/apiApp"]