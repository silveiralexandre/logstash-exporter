package requests

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/silveiralexandre/logstash-exporter/access"
)

// Client represents an HTTP client structure for performing requests
type Client struct {
	URI         string
	Retries     int
	OffsetLimit int
	httpClient  http.Client
	Timeout     time.Duration
	Proxy       string
}

// Setup will initiate the HTTP client with all required settings
func (r *Client) Setup(a access.Credential, uri string) (Client, error) {
	r.URI = uri
	r.Timeout = time.Duration(a.Timeout)
	r.Retries = a.Retries
	r.OffsetLimit = 50
	r.Proxy = os.Getenv("LSEXPORTER_PROXY")
	r.httpClient = http.Client{
		Transport: setTransport(r.Proxy),
		Timeout:   r.Timeout * time.Second,
	}
	return *r, nil
}

// Get executes a GET request to a target Endpoint and retries for a specific amount of times in case of errors
func (r *Client) Get(a access.Credential) ([]byte, error) {
	e := *new(error)
	b := []byte{}

	for r.Retries > 0 {
		b, e = doGet(a, r)
		if e != nil {
			r.Retries--
		} else {
			return b, nil
		}
	}
	return nil, e
}

func doGet(a access.Credential, r *Client) ([]byte, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*r.Timeout)
	defer cancel()

	req, err := http.NewRequest("GET", r.URI, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to wrap request with context: '%v'", err)
	}

	req.SetBasicAuth(a.Username, a.Password)
	req = req.WithContext(ctx)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read request body response: '%v'", err)
	}

	if len(b) == 0 || string(b) == "{}" {
		return nil, fmt.Errorf("Received empty response from target API")
	}
	return b, nil
}

func setTransport(proxy string) *http.Transport {
	if proxy != "" {
		proxyAddress, err := url.Parse(proxy)
		if err != nil {
			log.Fatal(fmt.Errorf("Could not parse proxy's URL from value provided: '%v'", err))
		}

		tr := &http.Transport{
			Proxy: http.ProxyURL(proxyAddress),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		return tr
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return tr
}
