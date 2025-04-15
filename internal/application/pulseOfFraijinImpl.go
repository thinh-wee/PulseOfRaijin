package application

import (
	"time"
)

func New(method, url string) PulseOfFraijinImpl {
	return &pulseOfFraiji{
		Method:             method,
		URL:                url,
		RequestsPerSecond:  5,
		MaxLifeTime:        30 * time.Second,
		InsecureSkipVerify: false,
		RequestTimeout:     30 * time.Second,
	}
}

type PulseOfFraijinImpl interface {
	// RunWithTPS runs the pulse of fraijin with the given TPS
	RunWithTPS(tps int) error
	// Start starts the pulse of fraijin
	Start() error

	// SetMaxLifeTime sets the max life time of the pulse of fraijin
	SetMaxLifeTime(maxLifeTime time.Duration) error
	// SetRequestsPerSecond sets the requests per second of the pulse of fraijin
	SetRequestsPerSecond(requestsPerSecond int) error
	// SetInsecureSkipVerify sets the insecure skip verify of the pulse of fraijin
	SetInsecureSkipVerify(insecureSkipVerify bool)
	// SetRequestTimeout sets the request timeout of the pulse of fraijin
	SetRequestTimeout(requestTimeout time.Duration) error
	// SetBody sets the body of the pulse of fraijin
	SetBody(body []byte) error
	// SetHeaders sets the headers of the pulse of fraijin
	SetHeaders(headers map[string]string) error
	// TODO: add a method to implement the pulse of fraijin

}
