FROM golang:latest

WORKDIR /app
COPY . .

RUN go build -o auth ./cmd/auth
CMD ["./auth"]
