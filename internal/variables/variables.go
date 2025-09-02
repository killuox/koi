package variables

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GetUserVariables() (map[string]interface{}, error) {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error finding home directory: %s", err)
	}

	// Build path ~/.koi/variables.json
	dirPath := filepath.Join(home, ".koi")
	filePath := filepath.Join(dirPath, "variables.json")

	// Ensure directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return nil, fmt.Errorf("error creating directory: %s", err)
		}
	}

	// If file doesn’t exist, create an empty JSON object
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.WriteFile(filePath, []byte("{}"), 0644)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %s", err)
		}
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	// Decode into a generic map
	var vars map[string]interface{}
	if err := json.Unmarshal(data, &vars); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}

	// Save back to file
	updated, _ := json.MarshalIndent(vars, "", "  ")
	if err := os.WriteFile(filePath, updated, 0644); err != nil {
		return nil, fmt.Errorf("error writing file: %s", err)
	}

	return vars, nil
}
