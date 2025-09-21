FROM golang:1.25-alpine AS builder

WORKDIR /oengus-timers
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /oengus-timers
COPY --from=builder /oengus-timers/main ./main
RUN chmod +x main

CMD ["./main"]
