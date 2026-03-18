#syntax=docker/dockerfile:1
FROM golang:1.26.1

WORKDIR /app

COPY go.mod go.sum ./

ENV RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

RUN go mod download

COPY gin.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend
EXPOSE 8080

# Run
CMD ["/backend"]
