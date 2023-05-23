FROM golang:1.20 AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN apt-get update && apt-get install -y libssl-dev

RUN go mod download

RUN go build  -o main cmd/main.go

FROM debian

WORKDIR /build

COPY --from=builder /build/main /build/main

EXPOSE 8080

CMD ["./main"]