package internal

import (
	tb "github.com/didip/tollbooth/v7"
	"io"
	"net/http"
	"time"
	"encoding/csv"
	"os"
	"encoding/json"
	"regexp"
)

func serve(config *Configuration, w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	switch r.Method {
	// reject all methods other than POST
	case http.MethodPost:
		contentType := r.Header.Get("Content-type")
		// make sure the data is JSON and we are in survey time.
		if contentType == "application/json" && now.Unix() < config.EndDate && now.Unix() > config.BeginDate {
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Unmarshal JSON to a string-string map.
			unmarshalledBody := make(map[string]string)
			err = json.Unmarshal([]byte(body), &unmarshalledBody)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Create slice array for CSV writing
			csvBuffer := make([]string, len(config.SurveyFields))

			// Validate JSON
			for index, field := range config.SurveyFields {
				// make sure all fields are filled out and make sure we don't have tampered data
				match, err := regexp.MatchString(`[,"\\]`, unmarshalledBody[field])
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

				if unmarshalledBody[field] == "" || match {
					w.WriteHeader(http.StatusBadRequest)
					return
				} else {
					csvBuffer[index] = unmarshalledBody[field]
				}
			}

			// create CSV handle
			csvHandle, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Write the CSV
			csvWriter := csv.NewWriter(csvHandle)
			defer csvWriter.Flush()
			csvWriter.Write(csvBuffer)

			// 202
			print("success")
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
	// See tollbooth docs for more info
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
