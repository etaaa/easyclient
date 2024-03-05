package easyclient

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Holds a *http.Client and the specified headers as those can't be saved on the client session directly.
type Client struct {
	client  *http.Client
	headers map[string]string
}

// Specifies default values when creating a new Client.
type ClientOptions struct {
	Cookies         ClientOptionsCookies // Custom cookies set on the client.
	FollowRedirects bool                 // If true, the client will follow redirects.
	Headers         map[string]string    // Custom headers set to the client.
	Proxy           string               // Custom proxy set to the client.
	Timeout         time.Duration        // Set custom request timeout. If nil, 30 * time.Second is used.
	Transport       http.RoundTripper    // Set custom transport type. If nil, http.DefaultTransport is used.
}

// RequestOptions specifies details when making a new request.
type RequestOptions struct {
	Body             io.Reader         // Body data to include in the request.
	Cookies          map[string]string // A map of cookies to set for the request.
	Headers          map[string]string // A map of headers to set for the request.
	Method           string            // The HTTP-Method.
	ReadResponseBody bool              // Whether the response.Body should be parsed or not.
	Proxy            string            // Custom proxy for the request.
	URL              string            // URL to perform the request on.
}

// Returns a new Client type.
func NewClient(options ClientOptions) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &Client{}, err
	}
	// Create new http.Client.
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: options.Timeout,
	}
	// Set default timeout to 30 seconds if not specified.
	if options.Timeout == 0 {
		httpClient.Timeout = 30 * time.Second
	}
	// Set custom transport if specified.
	if options.Transport != nil {
		httpClient.Transport = options.Transport
	}
	// Create client instance.
	c := &Client{
		client:  httpClient,
		headers: options.Headers,
	}
	// Set checkRedirect.
	if !options.FollowRedirects {
		c.SetFollowRedirects(false)
	}
	// Set proxy from string.
	if options.Proxy != "" {
		if err := c.SetProxy(options.Proxy); err != nil {
			return &Client{}, err
		}
	}
	// Set cookies.
	if err := c.SetCookies(options.Cookies.BaseURL, options.Cookies.Cookies); err != nil {
		return &Client{}, err
	}
	// Return client.
	return c, nil
}

// Executes the request.
func (c *Client) Do(options RequestOptions) (*http.Response, []byte, error) {
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
	// Set request proxy. Warning: This will overwrite the existing transport for the current request.
	tmpTransport := c.client.Transport
	if options.Proxy != "" {
		if err := c.SetProxy(options.Proxy); err != nil {
			return &http.Response{}, nil, err
		}
	}
	// Do the actual request.
	res, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, nil, err
	}
	// Close body on return.
	defer res.Body.Close()
	// Switch back to old client transport.
	c.client.Transport = tmpTransport
	// Parse the body.
	if options.ReadResponseBody {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return &http.Response{}, nil, err
		}
		return res, body, err
	}
	// Return reponse.
	return res, nil, err
}
