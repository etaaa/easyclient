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
func (c *Client) SetCookies(_url string, _cookies map[string]string) error {
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
func (c *Client) ClearCookies() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	c.client.Jar = jar
	return nil
}

// Returns all cookies set on the client.
func (c *Client) GetCookies(baseURL string) ([]*http.Cookie, error) {
	urlParsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return c.client.Jar.Cookies(urlParsed), nil
}
