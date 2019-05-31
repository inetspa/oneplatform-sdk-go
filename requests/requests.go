package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type Response struct {
	Code   int
	Status string
	Header http.Header
	Body   []byte
}

func Request(method string, url string, headers map[string]string, body io.Reader) (Response, error) {
	r := Response{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return r, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, err
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	r.Code = resp.StatusCode
	r.Status = resp.Status
	r.Body = buf.Bytes()
	r.Header = resp.Header
	_ = resp.Body.Close()

	return r, nil
}

func Get(url string, headers map[string]string) (Response, error) {
	return Request("GET", url, headers, nil)
}

func Post(url string, headers map[string]string, body io.Reader) (Response, error) {
	return Request("POST", url, headers, body)
}

func Put(url string, headers map[string]string, body io.Reader) (Response, error) {
	return Request("PUT", url, headers, body)
}

func Delete(url string, headers map[string]string, body io.Reader) (Response, error) {
	return Request("DELETE", url, headers, body)
}
