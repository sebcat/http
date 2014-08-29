
## http

A collection of tools for testing HTTP services as well as exposing some
Go functionality to the shell that makes life easier (e.g., the file server)

### Install

Make sure you have Go installed, e.g., by checking the Go version:

```
$ go version
go version go1.3.1 freebsd/amd64
```

If Go is not installed, install it using your OS package manager

If you havn't configured your GOPATH please do so:

```
$ mkdir $HOME/go-workspace
# put the envvars below in your shell rc file too
$ export GOPATH=$HOME/go-workspace
$ export PATH=$GOPATH/bin:$PATH
```

If Go is installed and your GOPATH is configured, get the package to your GOPATH

```
$ go get github.com/sebcat/http
```

And done!

```
$ http
Commands:
  http mwu
    Sample HTTP requests and perform the MW U test on two HTTP response-time groups
  http file-server
    Share a part of the local file system over HTTP
  http pong-server
    Start an HTTP server that responds with "pong\n"
  http get-urls
    Retrieve a list of HTTP resources and their status codes
  http stress-test
    Send HTTP requests at a specified rate and duration

```


## http mwu
Sample the response times for two requests (x,y) and calculate the p-value
for the [Mann-Whitney U test](http://en.wikipedia.org/wiki/Mann–Whitney_U_test). This can be used for mapping back end behavior 
e.g., finding side channels on the time spectrum, testing for blind SQL injections and
correlating changes of certain input parameters to an increased/decreased response time.

```
Usage:
  -request-timeout=20s: time-out value for a single request to complete
  -sample-size=20: number of requests per request type
  -throwaways=1: number of initially discarded request pairs
  -x-body="": request body for X
  -x-body-type="application/x-www-form-urlencoded ": request body type for X, if a request body is present
  -x-method="GET": HTTP request method for X
  -x-url="": URL for X
  -y-body="": request body for Y
  -y-body-type="application/x-www-form-urlencoded ": request body type for Y, if a request body is present
  -y-method="GET": HTTP request method for Y
  -y-url="": URL for Y
```

## http file-server

Share a part of the local file system over HTTP

```
Usage:
  -listen=":8080": listening directive
  -no-dir-list=false: disable directory listing
  -path="": path to HTTP root
```

## http pong-server
    
Start an HTTP server that responds with "pong\n"

```
Usage:
  -listen=":8080": listen directive
```

## http get-urls

Retrieve a list of HTTP resources and their status codes

```
Usage:
  -consume-body=false: consume http response body
  -http-timeout=20s: HTTP client timeout
  -method="HEAD": HTTP method
  -n-fetchers=20: number of concurrent HTTP fetchers
  -url-file="": file containing a newline separated list of URLs
```

## http stress-test

Send HTTP requests at a specified rate and duration

```
Usage:
  -body="": request body
  -body-type="": response body type
  -duration=3s: send duration
  -method="GET": HTTP method
  -rate=50: send rate (req/s)
  -timeout=20s: HTTP request timeout
  -url="": URL
```
