package dio

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var (
	envFlag        = flag.String("env", "development", "Set the application environment (reads environment from .env.{mode})")
	EnvFlag        string
	EnvFilePattern = ".env.%s" // default pattern
	printEnv       = true      // print environment
)

// InitEnvironment initealizes the environment after the caller parses flags
func InitEnvironment() {
	if EnvFlag == "" { // Don't overwrite if manually set
		EnvFlag = *envFlag
	}
	if err := loadEnvFile(EnvFlag); err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}
	printEnvMode(EnvFlag)
}

// InitEnvironmentWithEnv initializes the environment with the specified environment
func InitEnvironmentWithEnv(env string) {
	if env == "" {
		env = "development"
	}
	EnvFlag = env
	if err := loadEnvFile(env); err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}
	printEnvMode(env)
}

// SetEnvFilePattern sets the pattern for the environment file
func SetEnvFilePattern(pattern string) {
	EnvFilePattern = pattern
}

// SetPrintEnvMode sets whether to print the environment mode
func SetPrintEnvMode(enabled bool) {
	printEnv = enabled
}

// loadEnvFile loads the environment variables from the corresponding .env file
// fallsback to .env if .{mode}.env is not found
func loadEnvFile(mode string) error {
	file := fmt.Sprintf(EnvFilePattern, mode)
	// Try to load `.env.{mode}` first
	err := godotenv.Load(file)
	if err != nil && os.IsNotExist(err) {
		// Try fallback `.env`
		if fallbackErr := godotenv.Load(".env"); fallbackErr != nil && !os.IsNotExist(fallbackErr) {
			return fmt.Errorf("error loading fallback .env file: %w", fallbackErr)
		} else if fallbackErr == nil {
			// `.env` loaded successfully, so no error should be returned
			return nil
		}
		return fmt.Errorf("no environment file found for mode: %s, using fallback .env", mode)
	}
	// After successfully loading `.env.{mode}`, attempt to load `.env`
	_ = godotenv.Load(".env") // Load `.env` but ignore errors
	return nil
}

// printEnvMode prints the current application mode with appropriate color.
func printEnvMode(mode string) {
	if printEnv {
		c := color.New(color.FgBlue)
		if mode == "production" || mode == "prod" {
			c = color.New(color.FgRed)
		}
		c.Printf("Environment: %s\n", mode)
	}
}

// Env retrieves the value of the specified environment variable.
// If not set, it uses the fallback value if provided, or returns an error if not.
func Env(key string, fallback ...string) (string, error) {
	value := os.Getenv(key)
	if value != "" {
		return value, nil
	}
	if len(fallback) > 0 {
		return fallback[0], nil
	}
	return "", fmt.Errorf("environment variable %s is not set", key)
}

// MustEnv retrieves the value of the specified environment variable.
func MustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("[ERROR] Environment variable %s is required but not set", key)
	}
	return value
}

// Name returns the current environment name
func Name() string {
	return EnvFlag
}

// Dev checks if the current environment is development
func Dev() bool {
	return EnvFlag == "development" || EnvFlag == "dev"
}

// Prod checks if the current environment is production
func Prod() bool {
	return EnvFlag == "production" || EnvFlag == "prod"
}
