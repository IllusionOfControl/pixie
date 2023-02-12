# Build stage
FROM golang:1.17 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/app .

ENV TOKEN <insert-your-bot-token-here>

CMD ["./app"]