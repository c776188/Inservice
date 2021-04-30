package main

import (
	"Inservice/models"
	_ "Inservice/routers"
	"Inservice/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
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
	mapConfig, _ := config.NewConfig("ini", "conf/env.conf")
	searchUrl := mapConfig.String("SEARCH_URL")

	var result []models.IClass

	filename := "./data/" + time.Now().Format("2006-01-02") + ".json"
	// check path
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0777)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 爬蟲
		result = services.GetInitInservice(searchUrl).Class

		// 取得map距離
		// result = services.GetMapDuration(result)

		// write json
		file, err := json.MarshalIndent(result, "", " ")
		if err != nil {
			fmt.Println(err)
		}

		_ = ioutil.WriteFile(filename, file, 0777)
	}

	beego.Run()
}
