FROM golang:alpine

RUN apk update && apk add git

RUN go get github.com/MoSunDay/redix

EXPOSE 6380 7090

ENTRYPOINT ["redix"]

WORKDIR /root/