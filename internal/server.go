package internal

import (
	"fmt"
	"net/http"
)

func serveIndexPage(config *Configuration, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, config.IndexFile)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func serveTelemRequests(config *Configuration, w http.ResponseWriter, r *http.Request) {
}

func BeginListener(config *Configuration) error {
	fmt.Printf("%+v\n", config)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			serveIndexPage(config, w, r)
		},
	))

	mux.Handle("/submit", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			serveTelemRequests(config, w, r)
		},
	))

	var err error
	err = http.ListenAndServeTLS(":" + config.Port, config.PathToCertFile, config.PathToKeyFile, mux)
	if err != nil {
		return err
	}

	return nil
}
