# Dio

_You thought it was a README, but it was me, Dio._

![image](https://github.com/catgoose/dio/blob/dfa2af28f877bc5054b4b599ca7a0cb2a060d252/dio.png)

## About

Dio loads `.env.{mode}` environment files using [godotenv](https://github.com/joho/godotenv) as a dependency. Environment mode is set with commandline flags.

```bash
go run main.go -env production
```

## Installation

```bash
go get github.com/yourusername/dio
```

## Usage

1. Create your .env files, like:

```bash
.env.development
.env.josephjostar
.env.production
```

1. Import `Dio` and read environment. Dio will `log.Fatalf` if environment
   variable is not set

```go
package main

import (
 "fmt"
 "github.com/catgoose/dio"
)

func main() {
 // Dio loads the environment based on the flag passed (e.g., `-env=production,-env development`)
 // Default mode is `development`
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

Neovim btw
