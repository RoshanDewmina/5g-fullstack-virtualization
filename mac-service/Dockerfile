FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mac-service .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/mac-service /app/
EXPOSE 8082
CMD ["./mac-service"]
