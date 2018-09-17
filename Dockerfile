FROM golang:alpine

RUN apk add --no-cache --virtual git 
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/auyer/fastgate
WORKDIR /go/src/github.com/auyer/fastgate
RUN dep ensure
RUN  go install github.com/auyer/fastgate/ 

## If using config file:
#ADD ./config.json ./config.json
## If Using TLS
#ADD ./server.key ./server.key
#ADD ./server.pem ./server.pem

## Expose the desired ports
EXPOSE 8000 8443
## If using config file, add the -config flag with the location.
ENTRYPOINT /go/bin/fastgate
