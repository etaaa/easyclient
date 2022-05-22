package easyclient

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// Client holds a *http.Client type and the specified headers.
type client struct {
	client  *http.Client
	headers map[string]string
}

// ClientOptions specifies default values when creating a new Client.
type ClientOptions struct {
	CheckRedirect bool                 // If true, the client will follow redirects.
	Cookies       ClientOptionsCookies // Custom cookies set on the client.
	Headers       map[string]string    // Custom headers set to the client.
	Jar           http.CookieJar       // Set a custom cookiejar. If nil, a new one is created.
	Proxy         string               // Custom proxy set to the client.
	Timeout       time.Duration        // Set custom request timeout. If nil, 30 * time.Second is used.
	Transport     http.RoundTripper    // Set custom transport type. If nil, http.DefaultTransport is used.
}

type ClientOptionsCookies struct {
	BaseURL string            // The URL of the page to set the cookies on.
	Cookies map[string]string // A map of cookies to set.
}

// RequestOptions specifies details when making a new request.
type RequestOptions struct {
	Body      io.Reader         // Body data to include in the request.
	Cookies   map[string]string // A map of cookies to set for the request.
	Headers   map[string]string // A map of headers to set for the request.
	Method    string            // The HTTP-Method.
	ParseBody bool              // Whether the response.Body should be parsed or not.
	Proxy     string            // Custom proxy for the request.
	URL       string            // URL to perform the request on.
}

// Returns all cookies set on the client.
func (c *client) GetCookies(baseURL string) ([]*http.Cookie, error) {
	urlParsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return c.client.Jar.Cookies(urlParsed), nil
}

// Sets cookies on the client which will be reused on every request.
func (c *client) SetCookies(_url string, _cookies map[string]string) error {
	urlParsed, err := url.Parse(_url)
	if err != nil {
		return err
	}
	var cookies []*http.Cookie
	for key, value := range _cookies {
		cookies = append(cookies, &http.Cookie{
			Name:  key,
			Value: value,
		})
	}
	c.client.Jar.SetCookies(urlParsed, cookies)
	return nil
}

// Clears all cookies on the client.
func (c *client) ClearCookies() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	c.client.Jar = jar
	return nil
}

// Specify if the client should follow redirects.
func (c *client) SetCheckRedirect(shouldRedirect bool) {
	if shouldRedirect {
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		c.client.CheckRedirect = nil
	}
}

// Sets headers on the client which will be reused on every request.
func (c *client) SetHeaders(headers map[string]string) {
	c.headers = headers
}

// Clears all headers on the client.
func (c *client) ClearHeaders() {
	c.headers = nil
}

// Sets a proxy to the client.
func (c *client) SetProxy(proxy string) error {
	proxyParsed, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	c.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyParsed),
	}
	return nil
}

// Clears the proxy from the client by replacing the http.Transport type.
func (c *client) ClearProxy() {
	c.client.Transport = &http.Transport{}
}

// Sets the timeout duration for the client.
func (c *client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// Sets the transport for the client.
func (c *client) SetTransport(transport http.RoundTripper) {
	c.client.Transport = transport
}

// Executes the request.
func (c *client) Do(options RequestOptions) (*http.Response, []byte, error) {
	// Create the request.
	req, err := http.NewRequest(options.Method, options.URL, options.Body)
	if err != nil {
		return &http.Response{}, nil, err
	}
	// Set requets cookies.
	for key, value := range options.Cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}
	// Set client headers.
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	// Set request headers.
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}
	// Set request proxy.
	transportBefore := c.client.Transport
	if options.Proxy != "" {
		proxyParsed, err := url.Parse(options.Proxy)
		if err != nil {
			return &http.Response{}, nil, err
		}
		c.client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyParsed),
		}
	} else {
		c.client.Transport = &http.Transport{}
	}
	// Do the actual request.
	res, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, nil, err
	}
	// Close body on return.
	defer res.Body.Close()
	// Switch back to old client transport.
	c.client.Transport = transportBefore
	// Parse the body.
	if options.ParseBody {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &http.Response{}, nil, err
		}
		return res, body, err

	}
	// Return reponse.
	return res, nil, err
}

// Returns a new Client type
func NewClient(clientOptions ClientOptions) (*client, error) {
	// Create new http.Client
	httpClient := &http.Client{
		Jar:     clientOptions.Jar,
		Timeout: clientOptions.Timeout,
	}
	// Create cookiejar by default
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &client{}, err
	}
	if clientOptions.Jar == nil {
		httpClient.Jar = jar
	}

	// Set default timeout to 30 * time.Second if not specified
	if clientOptions.Timeout == 0 {
		httpClient.Timeout = 30 * time.Second
	}
	// Set custom transport if specified
	if clientOptions.Transport != nil {
		httpClient.Transport = clientOptions.Transport
	}
	// Create client instance
	c := &client{
		client:  httpClient,
		headers: clientOptions.Headers,
	}
	// Set checkRedirect
	if !clientOptions.CheckRedirect {
		c.SetCheckRedirect(false)
	}
	// Set proxy from string
	if clientOptions.Proxy != "" {
		if err := c.SetProxy(clientOptions.Proxy); err != nil {
			return &client{}, err
		}
	}
	// Set cookies
	if err := c.SetCookies(clientOptions.Cookies.BaseURL, clientOptions.Cookies.Cookies); err != nil {
		return &client{}, err
	}
	// Return client.
	return c, nil
}
