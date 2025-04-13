package application

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type pulseOfFraiji struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers map[string]string
	// RequestsPerSecond is the number of requests per second, default is 5
	RequestsPerSecond int

	// MaxLifeTime is the maximum life time of the pulse, default is 3 seconds
	MaxLifeTime time.Duration

	// InsecureSkipVerify is used to skip the certificate verification, default is false
	InsecureSkipVerify bool

	// RequestTimeout is the timeout of the request, default is 30 seconds
	RequestTimeout time.Duration
}

func (p *pulseOfFraiji) SetMaxLifeTime(maxLifeTime time.Duration) error {
	if maxLifeTime < 1*time.Second {
		return fmt.Errorf("max life time must be greater than 1 second")
	}
	p.MaxLifeTime = maxLifeTime
	return nil
}

func (p *pulseOfFraiji) SetRequestsPerSecond(requestsPerSecond int) error {
	if requestsPerSecond < 1 {
		return fmt.Errorf("requests per second must be greater than 0")
	}
	p.RequestsPerSecond = requestsPerSecond
	return nil
}

func (p *pulseOfFraiji) SetInsecureSkipVerify(insecureSkipVerify bool) {
	p.InsecureSkipVerify = insecureSkipVerify
}

func (p *pulseOfFraiji) SetRequestTimeout(requestTimeout time.Duration) error {
	if requestTimeout <= 0 {
		return fmt.Errorf("request timeout must be greater than 0")
	}
	p.RequestTimeout = requestTimeout
	return nil
}

func (p *pulseOfFraiji) SetBody(body io.Reader) error {
	if body == nil {
		return fmt.Errorf("body must not be nil or not set it")
	}
	p.Body = body
	return nil
}

func (p *pulseOfFraiji) SetHeaders(headers map[string]string) error {
	if headers == nil {
		return fmt.Errorf("headers must not be nil or not set it")
	}
	p.Headers = headers
	return nil
}

func (p *pulseOfFraiji) RunWithTPS(tps int) error {
	if tps < 1 {
		return fmt.Errorf("tps must be greater than 0")
	}
	if err := p.SetRequestsPerSecond(tps); err != nil {
		return err
	}
	return p.Start()
}

func (p *pulseOfFraiji) Start() error {

	var (
		sendRequestTimeLogs    []time.Time
		receiveRequestTimeLogs []time.Time

		wg sync.WaitGroup
	)

	request, err := http.NewRequest(p.Method, p.URL, p.Body)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: p.RequestTimeout,
	}

	if p.InsecureSkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// on 1 second, send n requests to the server;
	// then delay between requests in 1 second is 1/n second or 1000/n ms
	// so, if the number of requests per second is 100, the delay between requests is 10ms
	delayBetweenRequests := time.Second / time.Duration(p.RequestsPerSecond)

	// set the latest send request time for the first time; next time, it will set by the goroutine
	latestSendRequestTime := time.Now()

	// start logic
	startTime := time.Now()

	for time.Since(startTime) <= p.MaxLifeTime {
		// calculate the delay between requests
		differenceTime := time.Since(latestSendRequestTime)

		// if the delay is greater than the delay between requests, it means the requests are sent too fast
		if differenceTime > delayBetweenRequests {
			fmt.Println("delay is less than 0, it means the requests are sent too fast,"+
				"please check code logic to fix it - differenceTime:", differenceTime.String())
		}
		// wait for the delay
		t := <-time.After(delayBetweenRequests - differenceTime)

		// record the time
		latestSendRequestTime = t
		sendRequestTimeLogs = append(sendRequestTimeLogs, t)

		// start a goroutine to send request
		wg.Add(1)
		go func() {
			// TODO:  send request
			resp, err := httpClient.Do(request)
			if err != nil {
				fmt.Println("error sending request:", err)
				return
			}
			defer wg.Done()

			// record the receive request time
			receiveRequestTimeLogs = append(receiveRequestTimeLogs, time.Now())

			fmt.Println("[debug] response status:", resp.Status)
		}()
	}

	// wait for all requests to be sent and received
	wg.Wait()

	//

	println("--------------------------------")
	println("Configs variables:")
	println("Method:", p.Method)
	println("URL:", p.URL)
	println("RequestsPerSecond:", p.RequestsPerSecond)
	println("MaxLifeTime:", p.MaxLifeTime.String())
	println("InsecureSkipVerify:", p.InsecureSkipVerify)
	println("RequestTimeout:", p.RequestTimeout.String())
	println("--------------------------------")
	println("delayBetweenRequests:", delayBetweenRequests.String())
	println("--------------------------------")
	// print the send request time logs
	for i, t := range sendRequestTimeLogs {
		fmt.Printf("send request time (%d): %s\n", i+1, t.Format(time.RFC3339Nano))
	}

	for range receiveRequestTimeLogs {
		// TODO:  receive request
	}

	return nil
}
