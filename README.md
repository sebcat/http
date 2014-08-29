
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

```
$ http mwu -sample-size=12 \
        -x-url="http://wavsep.local/wavsep/active/SQL-Injection/SInjection-Detection-Evaluation-GET-200Error/Case19-InjectionInUpdate-NumericWithoutQuotes-CommandInjection-With200Errors.jsp?msgid=1%20xor%20(SELECT%20BENCHMARK(100000,%20MD5(0)))%20%20--%20" \
        -y-url="http://wavsep.local/wavsep/active/SQL-Injection/SInjection-Detection-Evaluation-GET-200Error/Case19-InjectionInUpdate-NumericWithoutQuotes-CommandInjection-With200Errors.jsp?msgid=1%20xor%20(SELECT%20BENCHMARK(1,%20MD5(0)))%20%20--%20"
        x                       y
        137.035266ms            66.845926ms
        136.625071ms            57.947573ms
        138.567795ms            57.925258ms
        136.0876ms              65.507054ms
        140.457938ms            60.340057ms
        136.803098ms            65.622593ms
        138.414367ms            65.411825ms
        136.523081ms            68.903624ms
        135.545243ms            56.019208ms
        145.325309ms            56.887758ms
        146.11732ms             74.017627ms
        139.016398ms            65.235549ms
p: 3.22564145623927e-05
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

```
$ cat urls.txt
http://www.google.com/
http://www.google.com.kh/
$ http get-urls -url-file=urls.txt
http://www.google.com.kh/ 200 OK 112.379692ms
http://www.google.com/ 200 OK 215.434101ms

```

## http stress-test

Send HTTP requests at a specified rate and duration

```
Usage:
  -body="": request body
  -body-type="": request body type
  -duration=3s: send duration
  -method="GET": HTTP method
  -rate=50: send rate (req/s)
  -timeout=20s: HTTP request timeout
  -url="": URL
```

```
$ http stress-test -rate=2 -duration=5.2s -url='http://www.google.com/'
total 10 (0 failed) acc time: 1.010624357s avg: 101.062435ms
```
