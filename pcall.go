package pcall

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/toretto460/pcall/batch"
)

type httpHandler func(w http.ResponseWriter, r *http.Request)

func errorHandler(w http.ResponseWriter) func(error) {
	return func(err error) {
		log.Print(err)

		switch e := err.(type) {
		case *batch.Error:
			w.WriteHeader(e.Code())
			w.Write([]byte(e.Error()))
		case error:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(e.Error()))
		}
	}
}

func NewHandler(allowedHosts []string) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		handledError := errorHandler(w)

		// Create batch requests from http request
		requests, err := batch.FromHTTPRequest(r)
		if err != nil {
			handledError(err)
			return
		}

		// Validate the batch against the whitelist
		whitelist := batch.WhitelistFromList(allowedHosts)
		err = requests.Validate(whitelist)
		if err != nil {
			handledError(err)
			return
		}

		res := requests.Send(http.Client{})

		content, jsonErr := json.Marshal(res)
		if jsonErr != nil {
			handledError(jsonErr)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}
