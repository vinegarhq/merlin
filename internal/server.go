package internal

import (
	"net/http"
)

func serve(config *Configuration, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, config.IndexFile)
	case http.MethodPost:
		contentType := r.Header.Get("Content-type")
		if contentType == "application/json" {
			//sanitize the shit HERE
			// WRITE THE CSV HERE
			// TODO: Figure out how to stop people from botting this shit
		} else {
			http.Error(w, "", http.StatusBadRequest)
		}
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func BeginListener(config *Configuration) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			serve(config, w, r)
		},
	))

	var err error
	err = http.ListenAndServeTLS(":"+config.Port, config.PathToCertFile, config.PathToKeyFile, mux)
	if err != nil {
		return err
	}

	return nil
}
