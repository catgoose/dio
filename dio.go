package dio

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

var envFlag = flag.String("env", "development", "Set the application environment (reads environment from .env.{mode})")

const (
	red   = "\033[31m"
	blue  = "\033[34m"
	reset = "\033[0m"
)

// InitEnvironment initializes the environment after the caller parses flags
func InitEnvironment() {
	// Use the parsed value of the environment flag
	EnvFlag := *envFlag

	if err := loadEnvFile(EnvFlag); err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	printEnvMode(EnvFlag)
}

// loadEnvFile loads the environment variables from the corresponding .env file
func loadEnvFile(mode string) error {
	file := fmt.Sprintf(".env.%s", mode)
	err := godotenv.Load(file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no environment file found for mode: %s", mode)
		}
		return fmt.Errorf("error loading environment file (%s): %w", file, err)
	}
	return nil
}

// printEnvMode prints the current application mode with appropriate color.
func printEnvMode(mode string) {
	pr := os.Getenv("DIO_PRINT_ENV")
	// Default to printing unless explicitly set to "false"
	if doPrint, err := strconv.ParseBool(pr); err != nil || doPrint {
		prodRegex := regexp.MustCompile(`^prod.*`)
		color := blue
		if prodRegex.MatchString(mode) {
			color = red
		}
		fmt.Printf("%sEnvironment: %s%s\n", color, mode, reset)
	}
}

// Env retrieves the value of the specified environment variable.
// If not set, it uses the fallback value if provided, or returns an error if not.
func Env(key string, fallback ...string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		if len(fallback) > 0 {
			return fallback[0], nil
		}
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}

// Name returns the current environment name
func Name() string {
	return *envFlag
}

// IsDev checks if the current environment is development
func Dev() bool {
	devRegex := regexp.MustCompile(`^dev.*`)
	return devRegex.MatchString(Name())
}

// IsProd checks if the current environment is production
func Prod() bool {
	prodRegex := regexp.MustCompile(`^prod.*`)
	return prodRegex.MatchString(Name())
}
