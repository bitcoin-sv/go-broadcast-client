// Package broadcast_client contains the client for the broadcast service.
package broadcast_client

// ArcClientConfig is used by [WithArc] to set up the connection between the broadcast client and Arc.
// The provided token will be used as the Authorization header.
type ArcClientConfig struct {
	APIUrl       string
	Token        string
	DeploymentID string
}

// GetApiUrl returns the API url.
func (c *ArcClientConfig) GetApiUrl() string {
	return c.APIUrl
}

// GetToken returns the token.
func (c *ArcClientConfig) GetToken() string {
	return c.Token
}

func (c *ArcClientConfig) GetDeploymentID() string {
	return c.DeploymentID
}
