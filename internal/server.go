package internal

import (
	"net/http"
)

func serveIndexPage(config *Configuration, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, config.indexFile)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func BeginListener(config *Configuration) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			serveIndexPage(config, w, r)
		},
	))

	var err error

	err = http.ListenAndServeTLS(":443", config.pathToCertFile, config.pathToKeyFile, mux)
	if err != nil {
		return err
	}

	return nil
}
