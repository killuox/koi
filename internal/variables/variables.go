package variables

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GetUserVariables() (map[string]any, error) {
	// Read the file
	data, err := getFile()
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
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filePath, updated, 0644); err != nil {
		return nil, fmt.Errorf("error writing file: %s", err)
	}

	return vars, nil
}

func SetUserVariable(key string, val any) error {
	// Read the file
	data, err := getFile()
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal existing variables
	var vars map[string]any
	if len(data) > 0 {
		if err := json.Unmarshal(data, &vars); err != nil {
			return fmt.Errorf("error unmarshaling file: %w", err)
		}
	} else {
		vars = make(map[string]any)
	}

	// Set or update the variable
	vars[key] = val

	// Marshal back to JSON
	updatedData, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling updated data: %w", err)
	}

	filePath, err := getFilePath()
	if err != nil {
		return err
	}

	// Write back to file
	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func getFile() ([]byte, error) {
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
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
	return data, nil
}

func getFilePath() (string, error) {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding home directory: %s", err)
	}

	// Build path ~/.koi/variables.json
	dirPath := filepath.Join(home, ".koi")
	filePath := filepath.Join(dirPath, "variables.json")

	// Ensure directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", fmt.Errorf("error creating directory: %s", err)
		}
	}

	return filePath, nil
}
