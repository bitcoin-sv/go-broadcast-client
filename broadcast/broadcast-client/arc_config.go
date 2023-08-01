package broadcast_client

type ArcClientConfig struct {
	APIUrl string
	Token  string
}

func (c *ArcClientConfig) GetApiUrl() string {
	return c.APIUrl
}

func (c *ArcClientConfig) GetToken() string {
	return c.Token
}
