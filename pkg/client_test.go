package easyclient

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	// Create a new client.
	client, _ := NewClient(ClientOptions{
		// Options specified here will persist for every request made with this client.
		Headers: map[string]string{
			"api-key": "123",
		},
	})
	// Execute a request with the created client.
	res, body, _ := client.Do(RequestOptions{
		// Options specified here will only be applied for the current request.
		Cookies: map[string]string{
			"foo": "bar",
		},
		Headers: map[string]string{
			"user-agent": "easyclient",
		},
		Method:    "GET",
		ParseBody: true,
		URL:       "https://httpbin.org/headers",
	})
	fmt.Print(res.StatusCode, string(body))
}
