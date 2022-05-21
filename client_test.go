package easyclient

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	// Client options
	client := NewClient()
	client.SetHeaders(map[string]string{
		"api-key": "123",
	})
	client.SetProxy("http://localhost:1234")
	// Request options
	res, body, err := client.Do(Options{
		Body: bytes.NewBuffer([]byte(`{"foo":"bar"}`)),
		Headers: map[string]string{
			"user-agent": "easy-client-go",
		},
		Method:    "POST",
		ParseBody: true,
		Proxy:     "",
		URL:       "https://httpbin.org/post",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.StatusCode, string(body))
}
