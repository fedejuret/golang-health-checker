FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/services /app/services

CMD ["./main"]
