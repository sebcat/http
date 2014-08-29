package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//   1) Read URL from file      (guUrlReader)
//   2) Fetch resource          (guResourceFetcher)
//   3) Consume resource        (guResourceConsumer)

// read URLs from file, pass to channel
func guUrlReader(file string) (<-chan string, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	ch := make(chan string)
	go func() {
		defer f.Close()
		for scanner.Scan() {
			ch <- scanner.Text()
		}

		close(ch)
	}()

	return ch, nil
}

// read urls from channel, fetch HTTP response and pass to channel
func guResourceFetcher(cli *http.Client, method string, urls <-chan string) <-chan *HttpPair {
	ch := make(chan *HttpPair)
	go func() {
		for url := range urls {
			r := &Request{Method: method, URL: url}
			req, err := r.ToHTTP() // Because of fields like User-Agent
			if err != nil {
				continue
			}

			start := time.Now()
			resp, err := cli.Do(req)
			t := time.Since(start)
			if err == nil {
				ch <- &HttpPair{Req: req, Resp: resp, Time: t}
			}
		}

		close(ch)
	}()

	return ch
}

// read HTTP request/response-pair from channel and do something with it (print it)
func guResourceConsumer(msgs <-chan *HttpPair, consumeBody bool) {
	for msg := range msgs {
		if consumeBody {
			ioutil.ReadAll(msg.Resp.Body)
		}

		fmt.Printf("%v %v %v\n", msg.Resp.StatusCode, msg.Time, msg.Req.URL)
		msg.Resp.Body.Close()
	}
}

func guMergeChans(ms []<-chan *HttpPair) <-chan *HttpPair {
	var wg sync.WaitGroup
	msgChan := make(chan *HttpPair)
	output := func(c <-chan *HttpPair) {
		for v := range c {
			msgChan <- v
		}

		wg.Done()
	}
	wg.Add(len(ms))
	for _, m := range ms {
		go output(m)
	}

	go func() {
		wg.Wait()
		close(msgChan)
	}()

	return msgChan
}

func guMain(args []string) error {
	var f flag.FlagSet
	var urlFile = f.String("url-file", "", "file containing a newline separated list of URLs")
	var nfetchers = f.Int("n-fetchers", 20, "number of concurrent HTTP fetchers")
	var httpTimeout = f.Duration("http-timeout", 20*time.Second, "HTTP client timeout")
	var consumeBody = f.Bool("consume-body", false, "consume http response body")
	var method = f.String("method", "HEAD", "HTTP method")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		} else {
			return err
		}
	}

	if len(*urlFile) == 0 {
		f.PrintDefaults()
		return nil
	}

	*method = strings.ToUpper(*method)
	cli := &http.Client{Timeout: *httpTimeout}
	urls, err := guUrlReader(*urlFile)
	if err != nil {
		return err
	}

	var ms []<-chan *HttpPair
	for i := 0; i < *nfetchers; i++ {
		m := guResourceFetcher(cli, *method, urls)
		ms = append(ms, m)
	}

	msgChan := guMergeChans(ms)
	guResourceConsumer(msgChan, *consumeBody)
	return nil
}
