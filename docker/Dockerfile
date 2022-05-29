FROM golang:latest

RUN mkdir -p /usr/src/lb
WORKDIR /usr/src/lb

COPY . /usr/src/lb
COPY configs/config.json /usr/src/lb/

RUN go build cmd/main.go

CMD [ "/usr/src/lb/main" ]
