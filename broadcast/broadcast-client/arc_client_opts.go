package broadcast_client

// ArcClientOptFunc defines an optional arguments that can be passed to the SubmitTransaction method.
type ArcClientOptFunc func(o *ArcClientOpts)

type ArcClientOpts struct {
	// XDeploymentID is the deployment id (id of instance using Arc API) that will be sent in the request header.
	XDeploymentID string
}

// GetApiUrl returns the API url.
func (c *ArcClientOpts) GetArcClientHeaders() map[string]string {
	headers := make(map[string]string)
	if c.XDeploymentID != "" {
		headers["X-Deployment-ID"] = c.XDeploymentID
	}

	return headers
}

// WithXDeploymentID is an option that allows you to set the deployment id (id of instance using Arc API) that will be sent in the request header.
func WithXDeploymentID(deploymentID string) ArcClientOptFunc {
	return func(o *ArcClientOpts) {
		o.XDeploymentID = deploymentID
	}
}
