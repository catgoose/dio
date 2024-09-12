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

var env string

const (
	red   = "\033[31m"
	blue  = "\033[34m"
	reset = "\033[0m"
)

func init() {
	flag.StringVar(&env, "env", "development", "Set the application environment (reads environment from .env.{mode})")
	flag.Parse()

	if err := loadEnvFile(env); err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	printEnvMode(env)
}

// loadEnvFile loads the environment variables from the corresponding .env file
func loadEnvFile(mode string) error {
	file := fmt.Sprintf(".env.%s", mode)
	err := godotenv.Load(file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("No environment file found for mode: %s", mode)
		}
		return fmt.Errorf("Error loading environment file (%s): %w", file, err)
	}
	return nil
}

// printEnvMode prints the current application mode with appropriate color.
func printEnvMode(mode string) {
	pr := Env("DIO_PRINT_MODE")
	// Default to printing unless explicitly set to "false"
	if doPrint, err := strconv.ParseBool(pr); err != nil || doPrint {
		prodRegex := regexp.MustCompile(`^prod.*`)
		color := blue
		if prodRegex.MatchString(mode) {
			color = red
		}
		fmt.Printf("%sApplication mode: %s%s\n", color, mode, reset)
	}
}

// Name returns the current environment name
func Name() string {
	return env
}

// Env retrieves the value of the specified environment variable
func Env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

// IsDev checks if the current environment is development
func Dev() bool {
	devRegex := regexp.MustCompile(`^dev.*`)
	return devRegex.MatchString(Name())
}

func Prod() bool {
	prodRegex := regexp.MustCompile(`^prod.*`)
	return prodRegex.MatchString(Name())
}
