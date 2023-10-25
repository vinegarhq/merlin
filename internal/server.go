package internal

import (
	"net/http"
	tb "github.com/didip/tollbooth/v7"
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

		} else {
			http.Error(w, "", http.StatusBadRequest)
		}
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func BeginListener(config *Configuration) error {
	mux := http.NewServeMux()

	tbLimiter := tb.NewLimiter(config.RateLimit, nil)
	tbLimiter.SetMethods([]string{"POST"})
	//NOTE, THE CONFIG BEHIND PROXIES IS DIFFERENT. YOU NEED TO ADD SETIPLOOKUPS
	mux.Handle("/", tb.LimitFuncHandler(tbLimiter, func(w http.ResponseWriter, r *http.Request) {
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
