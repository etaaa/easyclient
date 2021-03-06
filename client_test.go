package easyclient

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	//  Create client.
	client, _ := NewClient(ClientOptions{
		Headers: map[string]string{
			"api-key": "123",
		},
	})
	// Create request.
	res, body, _ := client.Do(RequestOptions{
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
	fmt.Println(res.StatusCode, string(body))
}
