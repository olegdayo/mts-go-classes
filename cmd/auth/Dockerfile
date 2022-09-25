FROM golang:latest

WORKDIR /app
COPY . .
RUN go build -o auth .
CMD ["./auth"]