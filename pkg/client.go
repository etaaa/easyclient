package easyclient

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Custom client type. Prevents all http.Client methods/properties to be accessible from outside.
type Client struct {
	client  *http.Client      // Private variable client as not all methods/properties should be accessible.
	headers map[string]string // Headers can't be stored on the *http.Client directly so we store them here.
}

// Specifies default values when creating a new Client.
type ClientOptions struct {
	Cookies         ClientOptionsCookies // Custom cookies which will be used on every request performed with this client for the given BaseURL.
	FollowRedirects bool                 // Whether the client should follow redirects or not.
	Headers         map[string]string    // Custom headers which will be used on every request performed with this client.
	Proxy           string               // Custom proxy which will be used on every request performed with this client.
	Timeout         time.Duration        // Timeout for when the request should be aborted if there is no response in that time. If nil, 30 * time.Second is used.
	Transport       http.RoundTripper    // Custom transport type. If nil, http.DefaultTransport is used.
}

// RequestOptions specifies details when making a new request.
type RequestOptions struct {
	Body             io.Reader         // Body to include in the request.
	Cookies          map[string]string // Custom cookies which will be used for this request.
	Headers          map[string]string // Custom headers which will be used for this request.
	Method           string            // The HTTP-Method. See https://pkg.go.dev/net/http#pkg-constants for all available methods.
	ReadResponseBody bool              // Whether the response.Body should be read or not.
	Proxy            string            // Custom proxy which will be used for this request.
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
