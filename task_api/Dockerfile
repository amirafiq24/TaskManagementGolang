FROM golang:1.27.1

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o task-api ./cmd/api

EXPOSE 8000

CMD ["./task-api"]