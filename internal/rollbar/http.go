package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// HTTPTimeout - timeout for HTTP client
	HTTPTimeout = time.Second * 60

	// TLSHandshakeTimeout - timeout for TLS handshake
	TLSHandshakeTimeout = time.Second * 30

	client = &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: TLSHandshakeTimeout,
		},
		Timeout: HTTPTimeout,
	}
)

// jcall - helper method for call json and parse object
func jcall(method, token string, url string, payload []byte, recv any) error {
	header := http.Header{
		"Accept":                 []string{"application/json"},
		"Content-Type":           []string{"application/json"},
		"X-Rollbar-Access-Token": []string{token},
	}

	r, err := call(method, url, header, payload)
	if err != nil {
		return err
	}
	defer r.Close()
	if recv == nil {
		return nil
	}
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(recv)
}

// call - HTTP call helper, returns the reader if success(200)
func call(method string, fullURL string, header http.Header, payload []byte) (io.ReadCloser, error) {
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method:        method,
		URL:           u,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(payload)),
		ContentLength: int64(len(payload)),
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logrus.Debugf("HTTP call failed - [%d]%s %s: %s", res.StatusCode, method, fullURL, res.Status)
		if res.Body != nil {
			bytes, err := io.ReadAll(res.Body)
			if err == nil {
				logrus.Debug("Body:")
				logrus.Debug(string(bytes))
			} else {
				logrus.Debugf("failed to read body - %v", err)
			}
			res.Body.Close()
		}
		return nil, fmt.Errorf("HTTP call failed - %d", res.StatusCode)
	}

	return res.Body, nil
}
