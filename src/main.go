package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	mux := http.NewServeMux()
	component := Index()
	mux.Handle("/", templ.Handler(component))

	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}
}
