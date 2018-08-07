
  [![FastGateLogo](https://raw.githubusercontent.com/auyer/FastGate/master/media/logo.png)](https://raw.githubusercontent.com/auyer/FastGate/master/media/logo.png)

  [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/auyer/FastGate)
  [![Release](https://img.shields.io/github/release/auyer/FastGate.svg)](https://github.com/auyer/fastgate/releases/latest) [![License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://github.com/auyer/FastGate/blob/master/LICENSE) [![Travis-CI](https://travis-ci.org/auyer/FastGate.svg?branch=master)](https://travis-ci.org/auyer/FastGate) [![Go Report Card](https://goreportcard.com/badge/github.com/auyer/FastGate?&fuckgithubcache=1)](https://goreportcard.com/report/github.com/auyer/FastGate)

## A Fast, light and Low Overhead API Gateway written in GO.

Fast, light and Low Overhead API Gateway written in GO
FastGate works by either proxying or redirecting traffic to the correct IP.

## Proxy VS Redirect

This Gateway was designed to work in two different modes. Here are a few differences, and how to chose them:

Proxy ( Enabled by Default) :

- Can reach every network acessible by the Gateway, independent on the Client and Service.
- All the data will be flowing trough this Gateway

Redirect ( Set `ProxyMode` to `false` in the Configuration ): 

- The connection with the Gateway will close right after the redirection, so the load will be minimum.
- This method works only if the Client and the Server could reach eachother in the first place.


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

## Example curl commands:

#### Registering 
```bash
  curl --request POST \
    --url http://localhost:8000/fastgate/ \
    --header 'content-type: application/json' \
    --data '{
      "address" : "http:/localhost:8080",
      "resource": "localapi"
      }
      '
```
#### Using route 

```bash
  curl --request GET \
    --url http://localhost:8000/api/localresource/ \
    --header 'x-fastgate-resource: localapi'
```
#### Geting List of Registered routes 

```bash
curl --request GET \
  --url http://localhost:8000/fastgate/ \
  --header 'content-type: application/json'
```


# TODO
- [x] Write a To-Do list
- [x] Create List All Routes for debugging
- [ ] Create a Proxy Option for outside Networks
- [ ] Benchmark comparing Redirect to Proxy
- [ ] Define scope ( How simple shoud we keep it ?)


