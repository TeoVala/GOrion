package envloader

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Config struct {
	AppName    string
	AppHost    string
	AppPort    string
	DBType     string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func LoadEnvVariables() *Config {
	f, err := os.Open(".env")
	if err != nil {
		// You might want to handle this error differently,
		// perhaps returning an error
		log.Printf("Could not open env file: %v", err)
		// Continue without a .env file, relying only on
		// already set environment variables
	} else {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			os.Setenv(key, val)
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading env file: %v", err)
		}
	}

	// Retrieve values from the OS environment after potential loading
	return &Config{
		AppName: os.Getenv("APP_NAME"),
		AppHost: os.Getenv("APP_HOST"),
		AppPort: os.Getenv("APP_PORT"),

		DBType:     os.Getenv("DB_TYPE"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}
