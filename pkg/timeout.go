package easyclient

import "time"

// Sets the timeout duration for the client.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}
