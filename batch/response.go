package batch

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Response represents the HTTP response
type Response struct {
	StatusCode int
	Content    string
	Headers    map[string]string
}

// NewResponse creates a new response
func NewResponse(code int, content []byte, headers map[string]string) Response {
	return Response{
		code,
		string(content),
		headers,
	}
}

// FromHTTPResponse creates a new Response from http.Response
func FromHTTPResponse(r *http.Response) Response {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	return Response{
		r.StatusCode,
		string(body),
		flatHeaders(r.Header),
	}
}

func flatHeaders(h http.Header) map[string]string {
	headers := make(map[string]string, 0)
	for name, value := range h {
		headers[name] = fmt.Sprintf("%v", value)
	}

	return headers
}
