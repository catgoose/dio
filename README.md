# Dio

<!--toc:start-->

- [Dio](#dio)
  - [About](#about)
  - [Installation](#installation)
  - [Usage](#usage) - [Environment](#environment)
  <!--toc:end-->

_You thought it was a README, but it was me, Dio._

![image](https://github.com/catgoose/screenshots/blob/b2cf4ef1674f99e894552af2c5cf654062ba4e37/dio/dio.png)

## About

Dio loads `.env.{mode}` environment files using [godotenv](https://github.com/joho/godotenv) as a dependency. Environment mode is set with commandline flags.

```bash
go run main.go -env production
```

## Installation

```bash
go get github.com/catgoose/dio
```

## Usage

1. Create your .env files, like:

```bash
.env.development
.env.josephjostar
.env.production
```

1. Import `Dio` and read environment.

```go
package main

import (
 "fmt"
 "github.com/catgoose/dio"
)

func main() {
 // Set your own flags
 // Parse the flags set by dio (-env)
 flag.Parse()
 // Dio loads the environment based on the flag passed (e.g., `-env=production,-env development`)
 // Default mode is `development`

 // Initialize environment after parsing flags
 dio.InitEnvironment()

 fmt.Println("Current environment:", dio.Name())

 // Access environment variables
 dbUser := dio.Env("DB_USER")
 fmt.Println("Database User:", dbUser)

 // Get environment name
 fmt.Printf("Environment %s", dio.Name())

 // Check environment
 if dio.Dev() {
  fmt.Println("We are in development mode.")
 }
 if dio.Prod() {
  fmt.Println("We are in production mode.")
 }
}
```

1. Set environment from commandline flag

```bash
go run main.go -env=production
go run main.go -env production
```

### Environment

To disable printing `Environment: ...` set environment variable `DIO_PRINT_ENV=false`
