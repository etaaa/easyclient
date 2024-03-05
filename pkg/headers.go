package easyclient

// Sets headers on the client which will be reused on every request.
func (c *Client) SetHeaders(headers map[string]string) {
	c.headers = headers
}

// Clears all headers on the client.
func (c *Client) ClearHeaders() {
	c.headers = nil
}
