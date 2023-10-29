package internal

import (
	tb "github.com/didip/tollbooth/v7"
	"io"
	"net/http"
	"time"
)

func serve(config *Configuration, w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodPost:
		contentType := r.Header.Get("Content-type")
		if contentType == "application/json" && now.Unix() < config.EndDate && now.Unix() > config.BeginDate {
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// WRITE THE CSV HERE
			print(string(body))
			w.WriteHeader(http.StatusAccepted)
			return
		} else {
			http.Error(w, "", http.StatusBadRequest)
			return
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
