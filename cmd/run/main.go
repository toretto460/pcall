package main

import (
	"github.com/toretto460/pcall"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	whitelist := []string{"reqres.in"}
	http.HandleFunc("/", pcall.NewHandler(whitelist))

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		log.Print(string(body))
		log.Print(r.Header)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
