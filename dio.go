// Package dio provides environment management utilities for Go applications.
// It supports loading environment variables from .env files based on different
// modes (development, production, etc.) and provides convenient functions
// for accessing environment variables with fallback values.
// The package requires the specific environment file (.env.{mode}) to exist
// and will return errors if not found, letting the client decide how to handle failures.
package dio

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// Error types for the dio package
var (
	// ErrEnvVarNotSet is returned when an environment variable is not set and no fallback is provided
	ErrEnvVarNotSet = errors.New("environment variable is not set")

	// ErrEnvFileNotFound is returned when no environment file can be found
	ErrEnvFileNotFound = errors.New("no environment file found")

	// ErrInvalidEnvMode is returned when an invalid environment mode is provided
	ErrInvalidEnvMode = errors.New("invalid environment mode")
)

var (
	envFlag        = flag.String("env", "development", "Set the application environment (reads environment from .env.{mode})")
	EnvFlag        string
	EnvFilePattern = ".env.%s" // default pattern
	printEnv       = true      // print environment
)

// InitEnvironment initializes the environment after the caller parses flags
func InitEnvironment() error {
	if EnvFlag == "" { // Don't overwrite if manually set
		EnvFlag = *envFlag
	}
	if err := loadEnvFile(EnvFlag); err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}
	printEnvMode(EnvFlag)
	return nil
}

// InitEnvironmentWithEnv initializes the environment with the specified environment
func InitEnvironmentWithEnv(env string) error {
	if env == "" {
		env = "development"
	}

	EnvFlag = env
	if err := loadEnvFile(env); err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}
	printEnvMode(env)
	return nil
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
// returns an error if the environment file is not found to prevent running in wrong environment
func loadEnvFile(mode string) error {
	if mode == "" {
		return fmt.Errorf("%w: mode cannot be empty", ErrInvalidEnvMode)
	}

	file := fmt.Sprintf(EnvFilePattern, mode)
	// Try to load `.env.{mode}`
	err := godotenv.Load(file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w for mode: %s (file: %s)", ErrEnvFileNotFound, mode, file)
		}
		return fmt.Errorf("error loading environment file %s: %w", file, err)
	}

	return nil
}

// printEnvMode prints the current application mode with appropriate color.
func printEnvMode(mode string) {
	if printEnv {
		c := color.New(color.FgBlue)
		if mode == "production" || mode == "prod" {
			c = color.New(color.FgRed)
		}
		_, err := c.Printf("Environment: %s\n", mode)
		if err != nil {
			// Fallback to standard output if color printing fails
			fmt.Printf("Environment: %s\n", mode)
		}
	}
}

// Env retrieves the value of the specified environment variable.
// If not set, it uses the fallback value if provided, or returns an error if not.
func Env(key string, fallback ...string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("%w: key cannot be empty", ErrEnvVarNotSet)
	}

	value := os.Getenv(key)
	if value != "" {
		return value, nil
	}
	if len(fallback) > 0 {
		return fallback[0], nil
	}
	return "", fmt.Errorf("%w: %s", ErrEnvVarNotSet, key)
}

// MustEnv retrieves the value of the specified environment variable.
// Returns an error if the variable is not set.
func MustEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%w: %s", ErrEnvVarNotSet, key)
	}
	return value, nil
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

// IsEnvVarNotSetError checks if the error is an environment variable not set error
func IsEnvVarNotSetError(err error) bool {
	return errors.Is(err, ErrEnvVarNotSet)
}

// IsEnvFileNotFoundError checks if the error is an environment file not found error
func IsEnvFileNotFoundError(err error) bool {
	return errors.Is(err, ErrEnvFileNotFound)
}

// IsInvalidEnvModeError checks if the error is an invalid environment mode error
func IsInvalidEnvModeError(err error) bool {
	return errors.Is(err, ErrInvalidEnvMode)
}

// EnvWithDefault retrieves the value of the specified environment variable with a default value
func EnvWithDefault(key, defaultValue string) string {
	value, err := Env(key, defaultValue)
	if err != nil {
		return defaultValue
	}
	return value
}
