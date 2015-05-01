package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"sort"
	"time"
)

type mwuSampleSettings struct {
	SampleSize     int
	RequestTimeout time.Duration
	NThrowaways    int
}

type mwuIndexTimePair struct {
	time  time.Duration
	index int
}

type mwuIndexTimePairs []mwuIndexTimePair

func (pairs mwuIndexTimePairs) Len() int {
	return len(pairs)
}

func (pairs mwuIndexTimePairs) Less(i, j int) bool {
	return pairs[i].time < pairs[j].time
}

func (pairs mwuIndexTimePairs) Swap(i, j int) {
	pairs[i], pairs[j] = pairs[j], pairs[i]
}

func mwuRankTime(xs, ys []time.Duration) (ranks []int) {
	pairs := make(mwuIndexTimePairs, len(xs)+len(ys))
	i := 0
	for ; i < len(xs); i++ {
		pairs[i].time = xs[i]
		pairs[i].index = i
	}

	for j := 0; i < len(pairs); i++ {
		pairs[i].time = ys[j]
		pairs[i].index = i
		j += 1
	}

	sort.Sort(pairs)
	ranks = make([]int, len(pairs))
	for i = 0; i < len(pairs); i++ {
		ranks[pairs[i].index] = i + 1
	}

	return ranks

}

// Mann-Whitney U test
// ties are not handled. This is normally not at problem for our
// purposes, but it should be noted
// Some reading:
// https://controls.engin.umich.edu/wiki/index.php/Basic_statistics:_mean,_median,_average,_standard_deviation,_z-scores,_and_p-value
// http://en.wikipedia.org/wiki/Mann%E2%80%93Whitney_U#Normal_approximation
func mwu(xs, ys []time.Duration) (p float64) {

	ranks := mwuRankTime(xs, ys)
	xranksum := 0
	for i := 0; i < len(xs); i++ {
		xranksum += ranks[i]
	}

	var (
		umin int
		u1   int = xranksum - (len(xs)*(len(xs)+1))/2
		u2   int = len(xs)*len(ys) - u1
	)

	if u1 < u2 {
		umin = u1
	} else {
		umin = u2
	}

	var (
		n1   int     = len(xs)
		n2   int     = len(ys)
		n1n2 float64 = float64(n1 * n2)
		eU   float64 = n1n2 / 2.0
		varU float64 = math.Sqrt(n1n2 * float64(n1+n2+1) / 12.0)
		z    float64 = (float64(umin) - eU) / varU
	)

	p = 1 + math.Erf(z/math.Sqrt(2))

	return p

}

func mwuSampleResponseTime(rt http.RoundTripper, r *Request) (t time.Duration,
	err error) {

	req, err := r.ToHTTP()
	if err != nil {
		return 0, err
	}

	start := time.Now()
	resp, err := rt.RoundTrip(req)
	if err != nil {
		return 0, err
	}

	t = time.Since(start)
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return t, nil
}

// fail on request error, KISS
func mwuSampleResponseTimes(xreq, yreq *Request, s *mwuSampleSettings) (
	xs, ys []time.Duration, err error) {
	xs = make([]time.Duration, s.SampleSize)
	ys = make([]time.Duration, s.SampleSize)
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: s.RequestTimeout,
		}).Dial}

	// first reqs will be outliers due to e.g., TCP handshake. discard them
	for i := 0; i < s.NThrowaways; i++ {
		if _, err = mwuSampleResponseTime(transport, xreq); err != nil {
			return nil, nil, err
		}

		if _, err = mwuSampleResponseTime(transport, yreq); err != nil {
			return nil, nil, err
		}
	}

	for i := 0; i < s.SampleSize; i++ {
		if xs[i], err = mwuSampleResponseTime(transport, xreq); err != nil {
			return nil, nil, err
		}

		if ys[i], err = mwuSampleResponseTime(transport, yreq); err != nil {
			return nil, nil, err
		}
	}

	return xs, ys, nil
}

func mwuMain(args []string) error {
	var settings mwuSampleSettings
	var xreq, yreq Request
	var f flag.FlagSet
	f.StringVar(&xreq.Method, "x-method", "GET", "HTTP request method for X")
	f.StringVar(&xreq.URL, "x-url", "", "URL for X")
	f.StringVar(&xreq.Body, "x-body", "", "request body for X")
	f.StringVar(&xreq.BodyType, "x-body-type", "application/x-www-form-urlencoded ",
		"request body type for X, if a request body is present")
	f.StringVar(&yreq.Method, "y-method", "GET", "HTTP request method for Y")
	f.StringVar(&yreq.URL, "y-url", "", "URL for Y")
	f.StringVar(&yreq.Body, "y-body", "", "request body for Y")
	f.StringVar(&yreq.BodyType, "y-body-type", "application/x-www-form-urlencoded ",
		"request body type for Y, if a request body is present")
	f.DurationVar(&settings.RequestTimeout, "request-timeout", 20*time.Second,
		"time-out value for a single request to complete")
	f.IntVar(&settings.SampleSize, "sample-size", 20,
		"number of requests per request type")
	f.IntVar(&settings.NThrowaways, "throwaways", 1,
		"number of initially discarded request pairs")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		} else {
			return err
		}
	}

	if len(xreq.URL) == 0 || len(yreq.URL) == 0 {
		f.PrintDefaults()
		return nil
	}

	if settings.SampleSize <= 0 {
		return errors.New("invalid sample size")
	}

	xsample, ysample, err := mwuSampleResponseTimes(&xreq, &yreq, &settings)
	if err != nil {
		return err
	}

	p := mwu(xsample, ysample)
	fmt.Printf("\tx\t\t\ty\n")
	for i := 0; i < settings.SampleSize; i++ {
		fmt.Printf("\t%v\t\t%v\n", xsample[i], ysample[i])
	}

	fmt.Printf("p: %v\n", p)
	return nil
}
