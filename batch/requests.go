package batch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	j "github.com/toretto460/pcall/batch/json"
)

// Requests is the array of requests
type Requests []*Request

// Request is the single request
type Request struct {
	Method  string
	URI     string
	URL     url.URL
	Content string
	Headers map[string]string
}

// Validate against the given hosts Whitelist
func (r *Requests) Validate(w Whitelist) error {
	for _, req := range *r {
		if !w.isWhitelisted(req.URL.Hostname()) {
			msg := fmt.Sprintf("Domain %s is not whitelisted.", req.URL.Hostname())
			return NewError(msg, 403)
		}
	}

	return nil
}

// UnmarshalJSON will decode the data into the Requests
func (r *Requests) UnmarshalJSON(data []byte) error {
	var jrqs j.Requests
	if err := json.Unmarshal(data, &jrqs); err != nil {
		return NewError("Cannot unmarshal requests: "+err.Error(), 422)
	}

	for _, jreq := range jrqs {
		req := &Request{
			Method:  jreq.Method,
			URI:     jreq.URI,
			Content: jreq.Content,
			Headers: jreq.Headers,
		}
		parsedURL, err := url.ParseRequestURI(req.URI)
		if err != nil {
			return NewError("Cannot parse uri "+req.URI, 422)
		}
		req.URL = *parsedURL

		*r = append(*r, req)
	}

	return nil
}

// Send will send all the requests in parallel manner
// and returns a slice of Responses
func (r *Requests) Send(client http.Client) []Response {
	responses := make(map[int]Response, len(*r))
	var wg sync.WaitGroup

	for id, req := range *r {
		wg.Add(1)
		go func(id int, req *Request) {
			defer wg.Done()
			forwardReq, err := http.NewRequest(
				req.Method,
				req.URL.String(),
				bytes.NewBuffer([]byte(req.Content)),
			)
			if err != nil {
				responses[id] = NewResponse(400, []byte(err.Error()), map[string]string{})
				return
			}

			for k, v := range req.Headers {
				forwardReq.Header.Set(k, v)
			}

			res, err := client.Do(forwardReq)

			if err != nil {
				responses[id] = NewResponse(402, []byte(err.Error()), map[string]string{})
				return
			}

			responses[id] = FromHTTPResponse(res)
		}(id, req)
	}

	wg.Wait()

	orderedResponses := make([]Response, len(*r))
	for k, r := range responses {
		orderedResponses[k] = r
	}

	return orderedResponses
}
