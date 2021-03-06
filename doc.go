// Package main controls all features of the FastGate API Gateway.
//
// To use this application, the user must send a POST request to /fastgate/ with the following body:
/*
{
  "address" 	: "https://yourEndpoint:8080"
  "resource"	: "resource-name"
}
*/
// This will create an entry in the database with the resource-name as a Key, and the address as the value.
// To access the desired route afterwards, add the X-fastgate-resource header with resource-name as its value.
//
// See this example below:
//
//
/*
$ curl --request POST   --url http://localhost:8000/fastgate/   --header 'content-type: application/json'   --data '{
  "address" : "http://localhost:8080",
  "resource"     : "hello"
}'

> HTTP/1.1 201 Created

$ curl --request GET \
  --url http://localhost:8000/hello \
  --header 'x-fastgate-resource: hello'

> GET /hello HTTP/1.1
> Host: localhost:8000
> x-fastgate-resource: hello
>
< HTTP/1.1 307 Temporary Redirect
< Location: http://localhost:8080/hello
$ Hello !


*/
// To help with usage and debugging, there is a route for listing all registered routes.
// This list will be returned when requesting a GET on /fastgate/
//
// See Example Below:
/*

$ curl --request GET   --url http://localhost:8000/fastgate/   --header 'content-type: application/json' -v
> GET /fastgate/ HTTP/1.1
...
< HTTP/1.1 202 Accepted
< Content-Type: application/json; charset=UTF-8
< Content-Length: 77
<
[
  {
    "address": "localapi",
    "resource": "http:/localhost:8080"
  }
*/
package main
