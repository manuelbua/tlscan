package tlscan

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func TlsConnect(host, port string, timeoutSeconds float64) (bool, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	dialer := &net.Dialer{
		Timeout: time.Duration(timeoutSeconds*1000) * time.Millisecond,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%s", host, port), conf)
	if err != nil {
		switch err.(type) {
		case tls.RecordHeaderError:
			return false, nil
		default:
			return false, err
		}
	}

	err = conn.Close()
	if err != nil {
		fmt.Printf("Error closing connection: %s", err)
	}
	return true, nil
}
