package easyclient

import "net/http"

// Sets the transport for the client.
func (c *Client) SetTransport(transport http.RoundTripper) {
	c.client.Transport = transport
}
