# easy-client

An easy to use Golang wrapper for the net/http package to perform requests.

## Features

* Auto body closing
* Auto body parsing
* Ordered headers
* Set headers for whole client
* Easy proxy setting for client
* Easy proxy setting for single request

## Usage

Install:
```bash
go get github.com/etaaa/easy-client
```

Usage:
```go
package main

import (
	"log"
	"github.com/etaaa/easy-client"
)

func main() {
	// Client options
	client := easyclient.NewClient()
	client.SetHeaders(map[string]string{
		"api-key": "123",
	})
	client.SetProxy("http://localhost:1234")
	// Request options
	res, body, err := client.Do(easyclient.Options{
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
```

## Questions
For any questions feel free to add and DM me on Discord (eta#0001).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
