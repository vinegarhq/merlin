package internal

import (
	"encoding/json"
	"encoding/csv"
	"os"
	"log"
)

type Configuration struct {
	Port           string   // Server port
	PathToCertFile string   // please use absolute paths
	PathToKeyFile  string   // please use absolute paths
	BeginDate      int64      // beginning date in epoch
	EndDate        int64      // end date in epoch
	OutputFile     string   // CSV file to record results to
	IndexFile      string   // User-facing index.html (optional)
	SurveyFields   []string // Yes, I am aware that this means all survey fields come out as strings, but this can be cleaned up in RStudio.
	RateLimit      float64      // The number of requests allowed per second
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
