package batch

import (
	"io/ioutil"
	"net/http"
)

// FromHTTPRequest parses the raw content and creates the Requests
func FromHTTPRequest(r *http.Request) (*Requests, error) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	var reqs Requests
	if err := reqs.UnmarshalJSON(body); err != nil {
		return nil, err
	}

	return &reqs, nil
}
