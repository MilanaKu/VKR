#1
FROM golang:1.24.4 AS builder
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

#2
FROM alpine:latest
WORKDIR /app 
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY marmelad.db .
EXPOSE 8080 
CMD ["/app/server"]