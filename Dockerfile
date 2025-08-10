FROM golang:1.24.5

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/app
RUN go build -o main main.go

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh


EXPOSE 8080
CMD ["/bin/sh", "-c", "/wait-for-it.sh postgres:5432 -- /wait-for-it.sh kafka:9092 -- /wait-for-it.sh redis:6379 -- ./main"]
