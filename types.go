package main

import (
	"net/http"
	"strings"
	"time"
)

type HttpPair struct {
	Req  *http.Request
	Resp *http.Response
	Time time.Duration // time from request to receiving the response header
}

// Representation of an HTTP request
// with JSON fields and a bit simpler than
// the net/http one
type Request struct {
	Method   string `json:"method"`
	URL      string `json:"url"`
	Body     string `json:"body,omitempty"`
	BodyType string `json:"body-type,omitempty"`
	Referer  string `json:"referer,omitempty"`
}

func (r *Request) ToHTTP() (httpReq *http.Request, err error) {
	method := strings.ToUpper(r.Method)
	if len(r.Body) > 0 {
		body := strings.NewReader(r.Body)
		httpReq, err = http.NewRequest(method, r.URL, body)
		if err != nil {
			return nil, err
		}

		if len(r.BodyType) > 0 {
			httpReq.Header.Set("Content-Type", r.BodyType)
		} else {
			httpReq.Header.Set("Content-Type",
				"application/x-www-form-urlencoded")
		}
	} else {
		httpReq, err = http.NewRequest(r.Method, r.URL, nil)
		if err != nil {
			return nil, err
		}
	}

	if len(r.Referer) > 0 {
		httpReq.Header.Set("Referer", r.Referer)
	}

	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (compatible; +github.com/sebcat/http)")
	return httpReq, nil
}
