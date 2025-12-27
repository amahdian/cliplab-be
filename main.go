package main

import (
	"fmt"
	"time"

	"github.com/amahdian/cliplab-be/cmd"
	"github.com/amahdian/cliplab-be/version"
)

func init() {
	// ensure server is always working in UTC timezone
	time.Local = time.UTC
}

// @title						AI Assistant
// @version					0.0.1
// @description				Swagger documentation for the AI Assistant's RESTful API.
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	fmt.Println("App version: ", version.Version())

	cmd.Execute()
}
