package easyclient

import (
	"net/http"
	"net/url"
)

// Sets a proxy to the client.
func (c *Client) SetProxy(proxy string) error {
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
func (c *Client) ClearProxy() {
	c.client.Transport = &http.Transport{}
}
