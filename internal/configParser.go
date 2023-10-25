package internal

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Port           string   // Server port
	PathToCertFile string   // please use absolute paths
	PathToKeyFile  string   // please use absolute paths
	BeginDate      int      // beginning date in epoch
	EndDate        int      // end date in epoch
	OutputFile     string   // CSV file to record results to
	IndexFile      string   // User-facing index.html (optional)
	SurveyFields   []string // Yes, I am aware that this means all survey fields come out as strings, but this can be cleaned up in RStudio.
	RateLimit      float64      // The number of requests allowed per second
	CSVHandle      *os.File // Please don't enter anything here!
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

	newConfig.CSVHandle, err = os.OpenFile(newConfig.OutputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return &Configuration{}, err
	}

	//TODO: add col names to CSV, depending on whether it is at the beginning of file.
	return newConfig, nil
}
