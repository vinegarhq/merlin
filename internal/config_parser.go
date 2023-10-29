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
	SurveyFields   []string `json:"surveyFields"`   // Yes, I am aware that this means all survey fields come out as strings, but this can be cleaned up in RStudio.
	RateLimit      float64  `json:"rateLimit"`      // The number of requests allowed per second
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

	// Write CSV header if file is blank.
	// CAUTION: DON'T MIX CSVS WITH DIFFERENT HEADERS
	info, err := os.Stat(newConfig.OutputFile)
	if err != nil {
		return &Configuration{}, err
	}

	if info.Size() == 0 {
		writer := csv.NewWriter(csvHandle)
		defer writer.Flush()
		writer.Write(newConfig.SurveyFields)
	} else {
		log.Println("caution: merlin will not write a header to an existing csv")
	}

	return newConfig, nil
}
