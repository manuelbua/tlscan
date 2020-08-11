package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

var (
	errNoHttp = errors.New("not an HTTP server")
	errCantConnect = errors.New("could not connect")
)

func TlsConnect(host, port string, timeoutSeconds float64) (bool, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	dialer := &net.Dialer{
		Timeout: time.Duration(timeoutSeconds*1000) * time.Millisecond,
	}

	hasTls := true
	conn, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%s", host, port), conf)
	if err != nil {
		switch err.(type) {
		case tls.RecordHeaderError:
			rhErr := err.(tls.RecordHeaderError)
			if strings.HasPrefix(string(rhErr.RecordHeader[:]), "HTTP/") {
				return false, nil
			}
			hasTls = false
		default:
			return false, err
		}
	}

	if conn == nil {
		return false, errCantConnect
	}

	// detect HTTP servers
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	if err != nil {
		// couldn't write request
		return hasTls, err
	}

	buf := make([]byte, 64)
	_, err = conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			// empty response, can't tell it's HTTP or not
			return hasTls, nil
		}

		// couldn't read response
		return hasTls, err
	}

	err = conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %s", err)
	}

	if len(buf) == 0 || strings.HasPrefix(string(buf), "HTTP/") {
		return hasTls, nil
	}

	return hasTls, errNoHttp
}
