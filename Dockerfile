FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go build -o kube-zodiakapp

EXPOSE 8080

CMD ./kube-zodiakapp
