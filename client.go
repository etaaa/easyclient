package easyclient

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client holds a *http.Client type and the specified headers
type Client struct {
	client  *http.Client
	headers map[string]string
}

// Options specifies details on the request
type Options struct {
	Body      io.Reader
	Headers   map[string]string
	Method    string
	ParseBody bool
	Proxy     string
	URL       string
}

// Specify if the http.Client type should follow redirects
// Default is set to false
func (client *Client) SetFollowRedirects(shouldRedirect bool) {
	if shouldRedirect {
		client.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		client.client.CheckRedirect = nil
	}
}

// Sets the headers on the Client type which will be reused on every request
func (client *Client) SetHeaders(headers map[string]string) {
	client.headers = headers
}

// Clears all headers on the Client type
func (client *Client) ClearHeaders() {
	client.headers = nil
}

// Sets a new proxy to the http.Client type
func (client *Client) SetProxy(proxy string) error {
	proxyParsed, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	client.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyParsed),
	}
	return nil
}

// Clears proxy from the http.Client type
func (client *Client) ClearProxy() {
	client.client.Transport = &http.Transport{}
}

// Sets the timeout duration for the http.Client type
// Default is set to 30 * time.Second
func (client *Client) SetTimeout(timeout time.Duration) {
	client.client.Timeout = timeout
}

// Executes the request
func (client *Client) Do(options Options) (*http.Response, []byte, error) {
	// Create the request
	req, err := http.NewRequest(options.Method, options.URL, options.Body)
	if err != nil {
		return &http.Response{}, nil, err
	}
	// Set headers from Client type
	for key, value := range client.headers {
		req.Header.Set(key, value)
	}
	// Set headers from request
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}
	// Set proxy
	transportBefore := client.client.Transport
	if options.Proxy != "" {
		proxyParsed, err := url.Parse(options.Proxy)
		if err != nil {
			return &http.Response{}, nil, err
		}
		client.client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyParsed),
		}
	} else {
		client.client.Transport = &http.Transport{}
	}
	// Do the actual request
	res, err := client.client.Do(req)
	if err != nil {
		return &http.Response{}, nil, err
	}
	// Close body on return
	defer res.Body.Close()
	// Switch back to old client transport
	client.client.Transport = transportBefore
	// Parse the body if wanted
	if options.ParseBody {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &http.Response{}, nil, err
		}
		return res, body, err

	}
	return res, nil, err
}

// Returns a new Client type
func NewClient() *Client {
	return &Client{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: 30 * time.Second,
		},
	}
}
