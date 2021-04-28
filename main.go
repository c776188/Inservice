package main

import (
	_ "Inservice/routers"

	"github.com/astaxie/beego"
)

// @title Inservice
// @version 1.0
// @description This is a sample server celler server.

// @host localhost:8080
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	beego.Run()
}
