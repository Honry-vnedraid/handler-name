FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o handler-service ./main.go

CMD ["./handler-service"]
