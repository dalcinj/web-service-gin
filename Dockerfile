FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-go-web-service-gin

EXPOSE 8080

CMD [ "/docker-go-web-service-gin" ]