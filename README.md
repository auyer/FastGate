
  [![FastGateLogo](https://raw.githubusercontent.com/auyer/FastGate/master/media/logo.png)](https://raw.githubusercontent.com/auyer/FastGate/master/media/logo.png)

  [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/auyer/FastGate)
  [![Release](https://img.shields.io/github/release/auyer/FastGate.svg)](https://github.com/auyer/fastgate/releases/latest) [![License](https://img.shields.io/badge/license-GPL3-brightgreen.svg)](https://github.com/auyer/FastGate/blob/master/LICENSE) [![Travis-CI](https://travis-ci.org/auyer/FastGate.svg?branch=master)](https://travis-ci.org/auyer/FastGate) [![Go Report Card](https://goreportcard.com/badge/github.com/auyer/FastGate?&fuckgithubcache=1)](https://goreportcard.com/report/github.com/auyer/FastGate)

## A Fast, light and Low Overhead API Gateway written in GO.

Fast, light and Low Overhead API Gateway written in GO
FastGate works by redirecting traffic to the correct IP. The connection to the Gateway closes just after the  redirect.

# Installation

To install fastgate, you can download the latest release binary from the [**Dowload page**](https://github.com/auyer/fastgate/releases/latest)
, or compile it from source with GO.

## Install Golang

If you need to install GO, please refer to the [golang.org](https://golang.org/dl/) Download Page, and follow instructions, or use a package manager (Most are very outdated). 

> For macOS users, I do recommend installing from homebrew. The mantainers are doing a amazing job keeping up with updates. Note that you still need to configure home path, but brew itself will teach you on how to do it.   Run : `brew install go`

## Fastgate Source instalation

```
go get github.com/auyer/FastGate
cd $GOPATH/src/github.com/auyer/FastGate
go install
```

# Deploy with Docker

By default, the Dockerfile picks the configuration file, TLS key and TLS cert from the same folder as the sourcecode.
```sh
  docker build -t fastgate .
  docker run -p YOUR_HTTP:8000 -p YOUR_HTTPS:8443 -d fastgate
```

# Usage
  ```
    fastgate -config ./path_to_config_file
  ```
  A sample to the configuration file can be found in [config.model.json](config.model.json)

#### To manually register (and test) FastGate, Send a POST request to `yourip:yourport/fastgate/` with a JSON like follows:
```
{
  "address" : "https://yourEndpoint:8080"
  "resource"     : "resource-name"
}
```
### Now send the desired request to `fastgate-ip:fastgate-port/your_resource` with the following header `X-fastgate-resource : resource-name`  and it should be working !



# TODO
- [ ] Write a To-Do list
