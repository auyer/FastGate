FROM golang:alpine

RUN apk add --no-cache --virtual git 
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/auyer/fastgate
WORKDIR /go/src/github.com/auyer/fastgate
RUN dep ensure
RUN  go install github.com/auyer/fastgate/ 

ADD ./config.json ./config.json
ADD ./server.key ./server.key
ADD ./server.pem ./server.pem

EXPOSE 8000 8443
ENTRYPOINT /go/bin/fastgate -config ./config.json
