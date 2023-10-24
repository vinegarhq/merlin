package internal

import (
	"encoding/csv"
	"os"
)

func OpenCSVHandle (config *Configuration) (*os.File, error){
	fileHandle, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// We set append-only to prevent accidental corruption of file.
	if err != nil {
		return nil, err
	}
	// REMEMBER TO CLOSE THE FILE HANDLE ON SHUTDOWN!!!
	return fileHandle, err
}
