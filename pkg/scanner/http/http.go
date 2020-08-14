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
	timeOut   time.Duration
	userAgent string
}

func NewScanner(timeoutSeconds float64) Scanner {
	return Scanner{
		timeOut:   time.Duration(timeoutSeconds*1000) * time.Millisecond,
		userAgent: fmt.Sprintf("tlscan/dudez"),
	}
}

func (s *Scanner) newClient(sni string) *http.Client {
	defaultDialFunc := &net.Dialer{
		Timeout:   s.timeOut,
		KeepAlive: time.Second,
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ServerName: sni}

	if net.ParseIP(sni) == nil {
		tlsConfig.ServerName = sni
	}

	// DialTLSContext could be an alternative
	// https://github.com/abursavich/dynamictls/blob/5d11b97955cdd8a1cb11f21e4012a26600cfa517/dynamictls.go#L270
	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   tlsConfig,
		DialContext:       defaultDialFunc.DialContext,
	}

	var client = &http.Client{
		Transport: tr,
		Timeout:   s.timeOut,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}

func (s *Scanner) Scan(ip, host, port string) (bool, error) {
	var url, urls string
	var hostname, hostHdr, sni string

	if len(ip) > 0 {
		hostname = ip
		sni = host
		hostHdr = host
		if port != "80" && port != "443" {
			hostHdr += ":" + port
		}
	} else {
		hostname = host
		sni = host
	}

	url = fmt.Sprintf("http://%s:%s", hostname, port)
	urls = fmt.Sprintf("https://%s:%s", hostname, port)
	client := s.newClient(sni)

	if isListening(client, s.userAgent, urls, hostHdr) {
		return true, nil
	} else if isListening(client, s.userAgent, url, hostHdr) {
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
