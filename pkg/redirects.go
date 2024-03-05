package easyclient

import "net/http"

// Specify if the client should follow redirects.
func (c *Client) SetRedirects(shouldRedirect bool) {
	if shouldRedirect {
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		c.client.CheckRedirect = nil
	}
}
