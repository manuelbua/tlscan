package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	errNoHttp      = errors.New("not an HTTP server")
)

type Scanner struct {
	client *http.Client
}

func New(timeoutSeconds float64) Scanner {
	timeout := time.Duration(timeoutSeconds*1000) * time.Millisecond

	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return Scanner{client: client}
}

func (h *Scanner) Scan(host, port string) (bool, error) {
	if isListening(h.client, fmt.Sprintf("https://%s:%s", host, port)) {
		return true, nil
	} else if isListening(h.client, fmt.Sprintf("http://%s:%s", host, port)) {
		return false, nil
	}
	return false, errNoHttp
}

func isListening(client *http.Client, url string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := client.Do(req)
	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	if err != nil {
		return false
	}

	return true
}
