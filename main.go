package main

import (
	_ "Inservice/routers"
	"Inservice/services"
	"os"
	"strconv"

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
	// clear data
	// files, _ := ioutil.ReadDir("./data")
	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// 	e := os.Remove("GeeksforGeeks.txt")
	// 	if e != nil {
	// 		fmt.Println("error delete data: ", e)
	// 	}
	// }

	// go job.CrawlerTaskListener()
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		beego.BConfig.Listen.HTTPPort = port
	}
	services.GetAndWriteInservice()
	beego.Run()
}
