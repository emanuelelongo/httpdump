# httpdump

Simple serivce that dumps HTTP requests.

## Description

This is a simple HTTP server that dumps the request body and headers to the console. It is useful for debugging and testing purposes.
It also send a response with the same body and headers to the client.

## Usage


### Compile and run

``` bash
go build -o httpdump main.go
./httpdump
```

### Specify a port

The default port is 8080. You can specify a different port by setting the `PORT` environment variable or by passing the `-port` flag.

Example:

``` bash
./httdump --port 8081

```


### Run with Docker

```
docker run -p 8080:8080 emanuelelongo/httpdump
```