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
	errRequestError = errors.New("request error")
)

type Scanner struct {
	client    *http.Client
	userAgent string
}

func NewScanner(timeoutSeconds float64) Scanner {
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

	var client = &http.Client{
		Transport: tr,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return Scanner{
		client:    client,
		userAgent: fmt.Sprintf("tlscan/dudez"),
	}
}

func (s *Scanner) Scan(ip, host, port string) (bool, error) {
	var url, urls string
	var hostname, hostHdr string

	if len(ip) > 0 {
		hostname = ip
		hostHdr = host
		if port != "80" && port != "443" {
			hostHdr += ":" + port
		}
	} else {
		hostname = host
	}

	url = fmt.Sprintf("http://%s:%s", hostname, port)
	urls = fmt.Sprintf("https://%s:%s", hostname, port)

	if isListening(s.client, s.userAgent, urls, hostHdr) {
		return true, nil
	} else if isListening(s.client, s.userAgent, url, hostHdr) {
		return false, nil
	}
	return false, errRequestError
}

func isListening(client *http.Client, ua, url, hostHdr string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	if len(hostHdr) > 0 {
		req.Host = hostHdr
	}

	req.Header.Add("Connection", "close")
	req.Close = true
	req.Header.Set("User-Agent", ua)

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
