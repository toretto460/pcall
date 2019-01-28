package json

// Requests is the array of requests
type Requests []*Request

// Request is the single request
type Request struct {
	Method  string            `json:"method,omitempty"`
	URI     string            `json:"uri,omitempty"`
	Content string            `json:"content,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}
