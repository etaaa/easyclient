package easyclient

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type ClientOptionsCookies struct {
	BaseURL string            // The URL of the page to set the cookies on.
	Cookies map[string]string // A map of cookies to set.
}

// Sets cookies on the client which will be reused on every request.
func (c *Client) SetCookies(URL string, cookieMap map[string]string) error {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	var cookies []*http.Cookie
	for name, value := range cookieMap {
		cookies = append(cookies, &http.Cookie{
			Name:  name,
			Value: value,
		})
	}
	c.client.Jar.SetCookies(parsedURL, cookies)
	return nil
}

// Clears all cookies on the client.
func (c *Client) ClearCookies() error {
	var err error
	c.client.Jar, err = cookiejar.New(nil)
	if err != nil {
		return err
	}
	return nil
}

// Returns all cookies set on the client.
func (c *Client) GetCookies(baseURL string) ([]*http.Cookie, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return c.client.Jar.Cookies(parsedURL), nil
}
