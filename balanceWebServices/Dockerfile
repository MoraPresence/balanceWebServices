FROM golang:latest

WORKDIR /app

COPY ./balanceWebServices/go.mod ./
COPY ./balanceWebServices/config.json ./

RUN go mod download

COPY ./balanceWebServices/*.go ./

RUN go build -o /balancer
CMD [ "/balancer" ]