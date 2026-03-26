package main

import (
	_ "raise-child/business"
	"raise-child/cmd"
)

// @title AgroTrust Server API
// @version 1.0
// @description API for AgroTrust Server
// @host locahost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
