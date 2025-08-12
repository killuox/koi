package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from the .env file.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, loading environment variables from OS")
	}
}

// GetString retrieves a string value from the environment with a fallback.
func GetString(key string, fallback string) (string, bool) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback, exists
	}
	return value, exists
}

// GetInt retrieves an integer value from the environment with a fallback.
func GetInt(key string, fallback int) (int, bool) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback, exists
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid integer for %s: %v. Using fallback %d\n", key, err, fallback)
		return fallback, exists
	}
	return intValue, exists
}

// GetBool retrieves a boolean value from the environment with a fallback.
func GetBool(key string, fallback bool) (bool, bool) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback, exists
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Invalid boolean for %s: %v. Using fallback %t\n", key, err, fallback)
		return fallback, exists
	}
	return boolValue, exists
}
