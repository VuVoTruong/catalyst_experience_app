package main

import (
	application "catalyst"
	"catalyst/config"
	"catalyst/docs"
	"fmt"
)

// @title Catalyst
// @version 1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath 7788/
func main() {
	cfg := config.NewConfig()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.ExposePort)

	application.Start(cfg)
}
