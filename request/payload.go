package request

type Payload struct {
	// URL to fetch
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}
