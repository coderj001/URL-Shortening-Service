FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o url-shortener

CMD ["./url-shortener"]

EXPOSE 3000

