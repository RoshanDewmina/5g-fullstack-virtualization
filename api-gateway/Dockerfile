FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-gateway .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/api-gateway /app/
EXPOSE 8080
CMD ["./api-gateway"]
