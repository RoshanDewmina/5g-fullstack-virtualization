FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o slice-manager .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/slice-manager /app/
EXPOSE 8084
CMD ["./slice-manager"]
