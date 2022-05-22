package easyclient

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	// Client options
	client, _ := NewClient()
	client.SetHeaders(map[string]string{
		"api-key": "123",
	})
	// Request options
	res, body, _ := client.DoRequest(Options{
		Cookies: map[string]string{
			"foo": "bar",
		},
		Headers: map[string]string{
			"user-agent": "easy-client",
		},
		Method:    "GET",
		ParseBody: true,
		URL:       "https://httpbin.org/headers",
	})
	fmt.Println(res.StatusCode, string(body))
}
