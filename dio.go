// Package dio provides environment management utilities for Go applications.
// It supports loading environment variables from .env files based on different
// modes (development, production, etc.) and provides convenient functions
// for accessing environment variables with fallback values.
// The package requires the specific environment file (.env.{mode}) to exist
// and will return errors if not found, letting the client decide how to handle failures.
package dio

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var (
	ErrEnvVarNotSet    = errors.New("environment variable is not set")
	ErrEnvFileNotFound = errors.New("no environment file found")
	ErrInvalidEnvMode  = errors.New("invalid environment mode")
)

var (
	envName       string
	canonicalMode string
)

const defaultEnv = "development"

// Options configures environment initialization. Nil means defaults
// (development mode, default file pattern, printing enabled).
type Options struct {
	// Env overrides the environment mode (e.g. "production", "uat").
	// Defaults to "development" if empty.
	Env string
	// FilePattern is the sprintf pattern for the env file (e.g. ".env.%s").
	// Defaults to ".env.%s" if empty.
	FilePattern string
	// Silent suppresses printing the environment mode at startup.
	// The zero value (false) means printing is enabled, matching nil opts behavior.
	Silent bool
}

const defaultFilePattern = ".env.%s"

func initEnv(env string, opts *Options) error {
	if env == "" {
		env = defaultEnv
	}
	envName = env
	pattern := defaultFilePattern
	printMode := true
	if opts != nil {
		if opts.FilePattern != "" {
			pattern = opts.FilePattern
		}
		printMode = !opts.Silent
	}
	if err := loadEnvFile(env, pattern); err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}
	canonicalMode = normalizeMode(env)
	printEnvMode(env, printMode)
	return nil
}

// InitEnvironment initializes the environment using opts.Env or the default
// environment ("development"). opts may be nil for defaults.
func InitEnvironment(opts *Options) error {
	env := defaultEnv
	if opts != nil && opts.Env != "" {
		env = opts.Env
	}
	return initEnv(env, opts)
}

// InitEnvironmentWithEnv initializes the environment with the given mode.
// opts may be nil for defaults.
func InitEnvironmentWithEnv(env string, opts *Options) error {
	return initEnv(env, opts)
}

func normalizeMode(mode string) string {
	switch strings.ToLower(mode) {
	case "dev":
		return "development"
	case "prod":
		return "production"
	case "production", "development", "uat":
		return strings.ToLower(mode)
	default:
		return strings.ToLower(mode)
	}
}

func loadEnvFile(mode, pattern string) error {
	if mode == "" {
		return fmt.Errorf("%w: mode cannot be empty", ErrInvalidEnvMode)
	}
	mode = normalizeMode(mode)
	file := fmt.Sprintf(pattern, mode)
	err := godotenv.Load(file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w for mode: %s (file: %s)", ErrEnvFileNotFound, mode, file)
		}
		return fmt.Errorf("failed to load environment file %s: %w", file, err)
	}

	return nil
}

func printEnvMode(mode string, print bool) {
	if !print {
		return
	}
	var attr color.Attribute
	switch canonicalMode {
	case "production":
		attr = color.FgRed
	case "uat":
		attr = color.FgYellow
	default:
		attr = color.FgBlue
	}
	_, err := color.New(attr).Printf("Environment: %s\n", mode)
	if err != nil {
		fmt.Printf("Environment: %s\n", mode)
	}
}

// Env retrieves the value of the environment variable named by key.
// If the variable is not set, it returns the first fallback value if provided,
// or an ErrEnvVarNotSet error if no fallback is given.
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

// RequiredEnv retrieves the value of the environment variable named by key.
// It returns an ErrEnvVarNotSet error if the variable is not set or empty.
func RequiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%w: %s", ErrEnvVarNotSet, key)
	}
	return value, nil
}

// Name returns the current environment mode name as provided during initialization.
func Name() string {
	return envName
}

// Dev reports whether the current environment is "development".
func Dev() bool {
	return canonicalMode == "development"
}

// Prod reports whether the current environment is "production".
func Prod() bool {
	return canonicalMode == "production"
}

// Uat reports whether the current environment is "uat".
func Uat() bool {
	return canonicalMode == "uat"
}

// IsEnvVarNotSetError reports whether err is or wraps ErrEnvVarNotSet.
func IsEnvVarNotSetError(err error) bool {
	return errors.Is(err, ErrEnvVarNotSet)
}

// IsEnvFileNotFoundError reports whether err is or wraps ErrEnvFileNotFound.
func IsEnvFileNotFoundError(err error) bool {
	return errors.Is(err, ErrEnvFileNotFound)
}

// IsInvalidEnvModeError reports whether err is or wraps ErrInvalidEnvMode.
func IsInvalidEnvModeError(err error) bool {
	return errors.Is(err, ErrInvalidEnvMode)
}

// EnvWithDefault retrieves the value of the environment variable named by key.
// If the variable is not set or empty, it returns defaultValue.
func EnvWithDefault(key, defaultValue string) string {
	if key == "" {
		return defaultValue
	}
	v, _ := Env(key, defaultValue)
	return v
}
