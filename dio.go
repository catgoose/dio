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
	envFlag       = flag.String("env", "development", "Set the application environment (reads environment from .env.{mode})")
	EnvFlag       string
	canonicalMode string
)

const defaultEnv = "development"

// Options configures init. Nil means defaults.
type Options struct {
	Env         string
	FilePattern string
	PrintMode   bool
}

const defaultFilePattern = ".env.%s"

func initEnv(env string, opts *Options) error {
	if env == "" {
		env = defaultEnv
	}
	EnvFlag = env
	pattern := defaultFilePattern
	printMode := true
	if opts != nil {
		if opts.FilePattern != "" {
			pattern = opts.FilePattern
		}
		printMode = opts.PrintMode
	}
	if err := loadEnvFile(env, pattern); err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}
	canonicalMode = normalizeMode(env)
	printEnvMode(env, printMode)
	return nil
}

// InitEnvironment initializes from the -env flag. opts may be nil for defaults.
func InitEnvironment(opts *Options) error {
	env := *envFlag
	if opts != nil && opts.Env != "" {
		env = opts.Env
	}
	return initEnv(env, opts)
}

// InitEnvironmentWithEnv initializes with the given env. opts may be nil for defaults.
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
		return mode
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

func RequiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%w: %s", ErrEnvVarNotSet, key)
	}
	return value, nil
}

func Name() string {
	return EnvFlag
}

func Dev() bool {
	return canonicalMode == "development"
}

func Prod() bool {
	return canonicalMode == "production"
}

func Uat() bool {
	return canonicalMode == "uat"
}

func IsEnvVarNotSetError(err error) bool {
	return errors.Is(err, ErrEnvVarNotSet)
}

func IsEnvFileNotFoundError(err error) bool {
	return errors.Is(err, ErrEnvFileNotFound)
}

func IsInvalidEnvModeError(err error) bool {
	return errors.Is(err, ErrInvalidEnvMode)
}

func EnvWithDefault(key, defaultValue string) string {
	v, err := Env(key, defaultValue)
	if err != nil {
		return defaultValue
	}
	return v
}
