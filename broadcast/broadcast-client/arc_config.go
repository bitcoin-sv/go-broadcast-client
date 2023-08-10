// Package broadcast_client contains the client for the broadcast service.
package broadcast_client

// ArcClientConfig is the reqired configuration for the ArcClient. Token will be used as the Authorization header.
type ArcClientConfig struct {
	APIUrl string
	Token  string
}

// GetApiUrl returns the API url.
func (c *ArcClientConfig) GetApiUrl() string {
	return c.APIUrl
}

// GetToken returns the token.
func (c *ArcClientConfig) GetToken() string {
	return c.Token
}
