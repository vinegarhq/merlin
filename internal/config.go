package internal

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Port           string   `json:"port"`           // Server port
	PathToCertFile string   `json:"pathToCertFile"` // please use absolute paths
	PathToKeyFile  string   `json:"pathToKeyFile"`  // please use absolute paths
	BeginDate      int64    `json:"beginDate"`      // beginning date in epoch
	EndDate        int64    `json:"endDate"`        // end date in epoch
	OutputFile     string   `json:"outputFile"`     // CSV file to record results to
	RateLimit      float64  `json:"rateLimit"`      // The number of requests allowed per second
	SurveyFields   []string `json:"surveyFields"`   // survey fields
}

func LoadConfiguration(pathToConfig string) (*Configuration, error) {
	newConfig := &Configuration{}

	cfg, err := os.ReadFile(pathToConfig)
	if err != nil {
		return &Configuration{}, err
	}

	err = json.Unmarshal(cfg, &newConfig)
	if err != nil {
		return &Configuration{}, err
	}

	csvHandle, err := os.OpenFile(newConfig.OutputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return &Configuration{}, err
	}

	// CAUTION: DON'T MIX CSVS WITH DIFFERENT HEADERS OR BAD THINGS WILL HAPPEN!!!
	info, err := os.Stat(newConfig.OutputFile)
	if err != nil {
		return &Configuration{}, err
	}

	// Write CSV headers if file is empty. Ideally, the user would specify a nonexistent file so that it can be created.
	if info.Size() == 0 {
		writer := csv.NewWriter(csvHandle)
		defer writer.Flush()
		err := writer.Write(newConfig.SurveyFields)
		if err != nil {
			return &Configuration{}, err
		}
	} else {
		log.Println("caution: merlin will not write a header to an existing csv")
	}

	return newConfig, nil
}
