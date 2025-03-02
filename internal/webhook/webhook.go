package webhook

// Webhook represents an incoming webhook payload with associated metadata.
type Webhook struct {
	Payload string `json:"payload"`
	Meta    Meta   `json:"meta"`
}

// Meta contains metadata information for a webhook.
type Meta struct {
	Client       string `json:"client"`
	Notification string `json:"notification"`
}
