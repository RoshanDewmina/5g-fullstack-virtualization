FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o pdcp-service .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/pdcp-service /app/
EXPOSE 8083
CMD ["./pdcp-service"]
