FROM golang:latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/cmd/consumer/main /app/cmd/consumer/main.go
EXPOSE 8080
CMD ["/app/cmd/consumer/main"]
