# Dio

[![Go Reference](https://pkg.go.dev/badge/github.com/catgoose/dio.svg)](https://pkg.go.dev/github.com/catgoose/dio)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<!--toc:start-->

- [Dio](#dio)
  - [About](#about)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [1. Create your .env files](#1-create-your-env-files)
    - [2. Basic Usage](#2-basic-usage)
    - [3. Run with environment flag](#3-run-with-environment-flag)
  - [API Reference](#api-reference)
    - [Environment Initialization](#environment-initialization)
      - [`InitEnvironment(opts *Options) error`](#initenvironmentopts-options-error)
      - [`InitEnvironmentWithEnv(env string, opts *Options) error`](#initenvironmentwithenvenv-string-opts-options-error)
      - [`Options`](#options)
    - [Environment Variable Access](#environment-variable-access)
      - [`Env(key string, fallback ...string) (string, error)`](#envkey-string-fallback-string-string-error)
      - [`RequiredEnv(key string) (string, error)`](#requiredenvkey-string-string-error)
      - [`EnvWithDefault(key, defaultValue string) string`](#envwithdefaultkey-defaultvalue-string-string)
    - [Environment Information](#environment-information)
      - [`Name() string`](#name-string)
      - [`Dev() bool`](#dev-bool)
      - [`Prod() bool`](#prod-bool)
      - [`Uat() bool`](#uat-bool)
    - [Configuration](#configuration)
  - [Error Handling](#error-handling)
    - [Error Types](#error-types)
    - [Error Handling Patterns](#error-handling-patterns)
      - [Fail-Fast (Recommended for Applications)](#fail-fast-recommended-for-applications)
      - [Graceful Handling (For Libraries)](#graceful-handling-for-libraries)
  - [Examples](#examples)
    - [Web Server Configuration](#web-server-configuration)
    - [Database Configuration](#database-configuration)
    - [Environment-Specific Behavior](#environment-specific-behavior)
  - [License](#license)
  <!--toc:end-->

_You thought it was a README, but it was me, Dio._

![image](https://github.com/catgoose/screenshots/blob/b2cf4ef1674f99e894552af2c5cf654062ba4e37/dio/dio.png)

## About

Dio is a Go package that provides environment management utilities for applications. It loads `.env.{mode}` environment files using [godotenv](https://github.com/joho/godotenv) and follows a fail-fast approach - the application will exit if the specified environment file doesn't exist, preventing accidental deployment to the wrong environment.

```bash
go run main.go -env production
```

## Installation

```bash
go get github.com/catgoose/dio
```

## Quick Start

### 1. Create your .env files

Create environment-specific files:

```bash
.env.development
.env.uat
.env.staging
.env.production
```

### 2. Basic Usage

```go
package main

import (
 "flag"
 "fmt"
 "log"

 "github.com/catgoose/dio"
)

func main() {
 flag.Parse()

 // Initialize environment - exits if .env.{mode} doesn't exist
 if err := dio.InitEnvironment(nil); err != nil {
  log.Fatalf("Failed to initialize environment: %v", err)
 }

 fmt.Printf("Running in %s mode\n", dio.Name())

 // Access environment variables
 port := dio.EnvWithDefault("PORT", "8080")
 dbHost := dio.EnvWithDefault("DB_HOST", "localhost")

 fmt.Printf("Server will run on port %s\n", port)
 fmt.Printf("Database host: %s\n", dbHost)
}
```

### 3. Run with environment flag

```bash
go run main.go -env=development
go run main.go -env=uat
go run main.go -env=production
go run main.go -env production
```

## API Reference

### Environment Initialization

#### `InitEnvironment(opts *Options) error`

Initializes the environment using the `-env` flag. Must be called after `flag.Parse()`. opts may be nil for defaults.

```go
flag.Parse()
if err := dio.InitEnvironment(nil); err != nil {
 log.Fatalf("Configuration error: %v", err)
}
```

#### `InitEnvironmentWithEnv(env string, opts *Options) error`

Initializes the environment with a specific environment name. opts may be nil for defaults.

```go
if err := dio.InitEnvironmentWithEnv("staging", nil); err != nil {
 log.Fatalf("Failed to load staging environment: %v", err)
}
```

#### `Options`

Options configures initialization. Nil means defaults: file pattern ".env.%s", PrintMode true.

- **Env** – override env (InitEnvironment only; when non-empty overrides the `-env` flag)
- **FilePattern** – env file pattern (e.g. `.env.%s`); empty means use default
- **PrintMode** – whether to print the environment mode on init

```go
if err := dio.InitEnvironmentWithEnv("production", &dio.Options{FilePattern: ".env.%s", PrintMode: false}); err != nil {
 log.Fatal(err)
}
```

### Environment Variable Access

#### `Env(key string, fallback ...string) (string, error)`

Retrieves an environment variable with optional fallback. Returns error if not set and no fallback provided.

```go
// With fallback
value, err := dio.Env("DB_HOST", "localhost")
if err != nil {
 // This won't happen because we provided a fallback
}

// Without fallback
apiKey, err := dio.Env("API_KEY")
if err != nil {
 log.Printf("API_KEY is required: %v", err)
}
```

#### `RequiredEnv(key string) (string, error)`

Retrieves a required environment variable. Returns error if not set.

```go
port, err := dio.RequiredEnv("PORT")
if err != nil {
 log.Fatalf("PORT is required: %v", err)
}
```

#### `EnvWithDefault(key, defaultValue string) string`

Convenience function that returns a default value if the environment variable is not set.

```go
port := dio.EnvWithDefault("PORT", "8080")
debug := dio.EnvWithDefault("DEBUG", "false")
```

### Environment Information

#### `Name() string`

Returns the current environment name.

```go
fmt.Printf("Current environment: %s\n", dio.Name())
```

#### `Dev() bool`

Checks if the current environment is development.

```go
if dio.Dev() {
 fmt.Println("Running in development mode")
}
```

#### `Prod() bool`

Checks if the current environment is production.

```go
if dio.Prod() {
 fmt.Println("Running in production mode")
}
```

#### `Uat() bool`

Checks if the current environment is UAT (User Acceptance Testing).

```go
if dio.Uat() {
 fmt.Println("Running in UAT mode")
}
```

### Configuration

Configuration is done via [Options](#options) when calling InitEnvironment or InitEnvironmentWithEnv.

## Error Handling

Dio provides consistent error handling with specific error types:

### Error Types

```go
// Check for specific error types
if dio.IsEnvVarNotSetError(err) {
 fmt.Println("Environment variable is not set")
}

if dio.IsEnvFileNotFoundError(err) {
 fmt.Println("Environment file not found")
}

if dio.IsInvalidEnvModeError(err) {
 fmt.Println("Invalid environment mode")
}
```

### Error Handling Patterns

#### Fail-Fast (Recommended for Applications)

```go
func main() {
 if err := dio.InitEnvironment(nil); err != nil {
  log.Fatalf("Configuration error: %v", err)
 }

 port, err := dio.RequiredEnv("PORT")
 if err != nil {
  log.Fatalf("Required env var missing: %v", err)
 }
}
```

#### Graceful Handling (For Libraries)

```go
func LoadConfig() (*Config, error) {
 if err := dio.InitEnvironment(nil); err != nil {
  return nil, fmt.Errorf("failed to initialize environment: %w", err)
 }

 port, err := dio.RequiredEnv("PORT")
 if err != nil {
  return nil, fmt.Errorf("failed to load config: %w", err)
 }

 return &Config{Port: port}, nil
}
```

## Examples

### Web Server Configuration

```go
package main

import (
 "flag"
 "fmt"
 "log"
 "net/http"

 "github.com/catgoose/dio"
)

func main() {
 flag.Parse()

 if err := dio.InitEnvironment(nil); err != nil {
  log.Fatalf("Failed to initialize environment: %v", err)
 }

 port := dio.EnvWithDefault("PORT", "8080")
 host := dio.EnvWithDefault("HOST", "0.0.0.0")

 addr := fmt.Sprintf("%s:%s", host, port)
 fmt.Printf("Starting server on %s in %s mode\n", addr, dio.Name())

 http.ListenAndServe(addr, nil)
}
```

### Database Configuration

```go
func LoadDBConfig() (*DBConfig, error) {
 config := &DBConfig{}

 // Required variables
 host, err := dio.RequiredEnv("DB_HOST")
 if err != nil {
  return nil, fmt.Errorf("database host required: %w", err)
 }
 config.Host = host

 // Optional with defaults
 config.Port = dio.EnvWithDefault("DB_PORT", "5432")
 config.Database = dio.EnvWithDefault("DB_NAME", "myapp")
 config.SSLMode = dio.EnvWithDefault("DB_SSL_MODE", "disable")

 return config, nil
}
```

### Environment-Specific Behavior

```go
func setupLogging() {
 if dio.Dev() {
  // Development: verbose logging
  log.SetLevel(log.DebugLevel)
 } else if dio.Uat() {
  // UAT: detailed logging for testing
  log.SetLevel(log.InfoLevel)
 } else if dio.Prod() {
  // Production: minimal logging
  log.SetLevel(log.WarnLevel)
 } else {
  // Other environments: info level
  log.SetLevel(log.InfoLevel)
 }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
