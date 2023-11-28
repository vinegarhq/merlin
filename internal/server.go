// test with
// curl https://127.0.0.1:7000 -k -X POST -H "Content-Type: application/json" -d @testJson.json
package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	tb "github.com/didip/tollbooth/v7"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

func serve(config *Configuration, w http.ResponseWriter, r *http.Request, regexpPointer *regexp.Regexp) {
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
			unmarshalledBody := make(map[string]interface{})
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
				singleField := fmt.Sprint(unmarshalledBody[field])
				match := regexpPointer.MatchString(singleField)
				if singleField == "" || match {
					w.WriteHeader(http.StatusBadRequest)
					return
				} else {
					csvBuffer[index] = singleField
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
			err = csvWriter.Write(csvBuffer)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// 202
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

	regexpPointer, err := regexp.Compile(`[^A-Za-z0-9./()]+`)
	if err != nil {
		return err
	}

	//NOTE, THE CONFIG BEHIND PROXIES IS DIFFERENT. YOU NEED TO ADD SETIPLOOKUPS
	// See tollbooth docs for more info
	mux.Handle("/", tb.LimitFuncHandler(tbLimiter, func(w http.ResponseWriter, r *http.Request) {
		serve(config, w, r, regexpPointer)
	},
	))

	err = http.ListenAndServeTLS(":"+config.Port, config.PathToCertFile, config.PathToKeyFile, mux)
	if err != nil {
		return err
	}

	return nil
}
