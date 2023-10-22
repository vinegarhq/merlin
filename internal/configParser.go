package internal

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	pathToCertFile string   // please use absolute paths
	pathToKeyFile  string   // please use absolute paths
	beginDate      int      // beginning date in epoch
	endDate        int      // end date in epoch
	outputFile     string   // CSV file to record results to
	indexFile      string   // User-facing index.html (optional)
	surveyFields   []string // Yes, I am aware that this means all survey fields come out as strings, but this can be cleaned up in RStudio.
}

func LoadConfiguration(pathToConfig string) (*Configuration, error) {
	newConfig := &Configuration{}
	cfg, err := os.ReadFile(pathToConfig)
	if err != nil {
		return newConfig, err
	}

	err = json.Unmarshal([]byte(cfg), &newConfig)
	if err == nil {
		return newConfig, err
	}
	return newConfig, nil
}
