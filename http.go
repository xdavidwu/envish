package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type progressLogger struct {
	of             io.ReadCloser
	expected, read int64
	logger         logr.Logger
	quit           chan struct{}
}

var _ io.ReadCloser = &progressLogger{}

func (p *progressLogger) Start() {
	timer := time.Tick(time.Second)
	p.quit = make(chan struct{})

	go func() {
		for {
			select {
			case <-timer:
				if p.expected != 0 {
					p.logger.Info(
						fmt.Sprintf("downloading %f%%", float64(p.read)/float64(p.expected)*100),
						"downloaded", p.read,
						"total", p.expected,
					)
				} else {
					p.logger.Info("downloading", "downloaded", p.read)
				}
			case <-p.quit:
				return
			}
		}
	}()
}

func (p *progressLogger) Stop() {
	p.quit <- struct{}{}
}

func (p *progressLogger) Read(b []byte) (int, error) {
	n, err := p.of.Read(b)
	p.read += int64(n)
	return n, err
}

func (p *progressLogger) Close() error {
	p.Stop()
	return p.of.Close()
}

type loggingRoundTripper struct {
	http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	before := time.Now()
	logger := log.Log.V(1).WithValues(
		"method", r.Method,
		"url", r.URL.String(),
		"at", before,
	)

	logger.Info("request")
	res, err := l.RoundTripper.RoundTrip(r)

	duration := time.Since(before)

	if err != nil {
		logger.Error(err, "request error")
	} else {
		logger.Info("received header", "code", res.StatusCode, "header_duration", duration)
	}

	if res.Header.Get("Content-Disposition") != "" {
		l := &progressLogger{of: res.Body, expected: res.ContentLength, logger: logger}
		l.Start()
		res.Body = l
	}

	return res, err
}

func requestLog(r http.RoundTripper) http.RoundTripper {
	return &loggingRoundTripper{r}
}
