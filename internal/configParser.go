package internal

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	PathToCertFile string   // please use absolute paths
	PathToKeyFile  string   // please use absolute paths
	BeginDate      int      // beginning date in epoch
	EndDate        int      // end date in epoch
	OutputFile     string   // CSV file to record results to
	IndexFile      string   // User-facing index.html (optional)
	SurveyFields   []string // Yes, I am aware that this means all survey fields come out as strings, but this can be cleaned up in RStudio.
}

func LoadConfiguration(pathToConfig string) (*Configuration, error) {
	newConfig := &Configuration{}
	cfg, err := os.ReadFile(pathToConfig)
	if err != nil {
		return newConfig, err
	}

	err = json.Unmarshal(cfg, &newConfig)
	if err != nil {
		return newConfig, err
	}

	return newConfig, nil
}
