# base go image
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o api-service ./cmd/api

RUN chmod +x /app/api-service

# Build a tiny docker image
FROM scratch

COPY --from=builder /app/api-service .

CMD [ "./api-service" ]