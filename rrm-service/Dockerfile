FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rrm-service .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/rrm-service /app/
EXPOSE 8081
CMD ["./rrm-service"]
