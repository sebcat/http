package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	STATUS_SUCCESS = iota
	STATUS_FAILED
)

type stReqStat struct {
	status int
	time   time.Duration
}

type stSenderStat struct {
	nsucceded int
	nfailed   int
	time      time.Duration
}

func stSendHttpRequest(cli *http.Client, req *Request, statusChan chan stReqStat) {
	httpReq, err := req.ToHTTP()
	status := stReqStat{status: STATUS_FAILED}
	if err == nil {
		startTime := time.Now()
		resp, err := cli.Do(httpReq)
		if err == nil {
			ioutil.ReadAll(resp.Body) // read entire body
			resp.Body.Close()
			status.time = time.Since(startTime)
			if resp.StatusCode == 200 {
				status.status = STATUS_SUCCESS
			}
		} else {
			status.time = time.Since(startTime)
		}
	}

	statusChan <- status
}

func stSendHttpRequests(req *Request, sendRate int, duration, timeout time.Duration) *stSenderStat {

	client := &http.Client{Timeout: timeout}
	ticker := time.NewTicker(time.Second / time.Duration(sendRate))
	httpStatusChan := make(chan stReqStat)
	doneSendChan := time.After(duration)
	var sstat stSenderStat
	waitGroup := &sync.WaitGroup{}
	go func() {
		for stReqStat := range httpStatusChan {
			if stReqStat.status == STATUS_SUCCESS {
				sstat.nsucceded += 1
			} else {
				sstat.nfailed += 1
			}

			sstat.time += stReqStat.time
			waitGroup.Done()
		}
	}()

	for {
		select {
		case <-ticker.C:
			waitGroup.Add(1)
			go stSendHttpRequest(client, req, httpStatusChan)
		case <-doneSendChan:
			ticker.Stop()
			waitGroup.Wait()
			close(httpStatusChan)
			return &sstat
		}
	}
}

func stMain(args []string) error {
	var req Request
	var f flag.FlagSet
	sendRate := f.Int("rate", 50, "send rate (req/s)")
	duration := f.Duration("duration", 3*time.Second, "send duration")
	timeout := f.Duration("timeout", 20*time.Second, "HTTP request timeout")
	f.StringVar(&req.Method, "method", "GET", "HTTP method")
	f.StringVar(&req.URL, "url", "", "URL")
	f.StringVar(&req.Body, "body", "", "request body")
	f.StringVar(&req.BodyType, "body-type", "", "request body type")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		} else {
			return err
		}
	}

	if len(req.URL) == 0 {
		f.PrintDefaults()
		return nil
	}

	res := stSendHttpRequests(&req, *sendRate, *duration, *timeout)
	tot := res.nsucceded + res.nfailed
	avg := res.time / time.Duration(tot)
	fmt.Printf("total %v (%v failed) acc time: %v avg: %v\n", tot, res.nfailed,
		res.time, avg)
	return nil
}
