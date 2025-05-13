package main

import "simple-finance/internal/app"
import _ "net/http"

//go:generate swag init -o=./swagger --parseDependency --parseDepth=1
// рабочая команда для запуска из корня проекта
//swag init -d "./" -g "cmd/web/main.go" --parseDependency --parseInternal --parseDepth 1

// Package docs Simple Finance API
//
// # Documentation for Simple Finance API
//
// @title           Simple Finance API
// @version         1.0
// @description     This is a simple finance management API
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8000
// @BasePath  /
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description                 Type "Bearer" followed by a space and JWT token
func main() {
	app.Run()
}
