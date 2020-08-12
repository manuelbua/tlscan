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
	errNoHttp = errors.New("not an HTTP server")
)

type Scanner struct {
	timeout time.Duration
}

func (s *Scanner) newClient(sni string) *http.Client {
	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true, ServerName: sni},
		DialContext: (&net.Dialer{
			Timeout:   s.timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}

	return &http.Client{
		Transport: tr,
		Timeout:   s.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func New(timeoutSeconds float64) Scanner {
	timeout := time.Duration(timeoutSeconds*1000) * time.Millisecond
	return Scanner{timeout: timeout}
}

func (s *Scanner) Scan(ip, host, port string) (bool, error) {
	var client *http.Client
	var url, urls string
	if len(ip) > 0 {
		client = s.newClient(host)
		url = fmt.Sprintf("http://%s:%s", ip, port)
		urls = fmt.Sprintf("https://%s:%s", ip, port)
	} else {
		client = s.newClient("")
		url = fmt.Sprintf("http://%s:%s", host, port)
		urls = fmt.Sprintf("https://%s:%s", host, port)
	}

	if isListening(client, urls) {
		return true, nil
	} else if isListening(client, url) {
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
