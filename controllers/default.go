package controllers

import (
	"Inservice/models"
	"Inservice/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/config"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// index 頁面
func (c *MainController) Get() {
	c.TplName = "index.tpl"
}

// @Title Class Info
// @Description 列表上課內容
// @
// @Success  200  object  models.IClass  "上課資訊"
// @
// @Resource 關於上課內容
// @Router / [post]
func (c *MainController) Post() {
	var result []models.IClass

	// check path
	filename := "./data/" + time.Now().Format("2006-01-02") + ".json"
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0777)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		result = services.GetAndWriteInservice()
	} else {
		// read json
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
		}

		_ = json.Unmarshal([]byte(file), &result)
	}

	c.Data["json"] = &result
	c.ServeJSON()
}
