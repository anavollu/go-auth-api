FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .
RUN go build -o app

FROM alpine:latest as release
WORKDIR /app
RUN mkdir html
COPY --from=builder /app/app .
COPY --from=builder /app/html ./html
CMD ["./app"]