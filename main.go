package main

import (
	_ "raise-child/business"
	"raise-child/cmd"
)

// @title AgroTrust Server API
// @version 1.0
// @description API for AgroTrust Server
// @host https://agrotrust-fkhwgwh5ecgwhrbf.indonesiacentral-01.azurewebsites.net/
// @BasePath /
// @schemes https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
