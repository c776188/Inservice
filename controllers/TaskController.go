package controllers

import (
	"Inservice/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/config"
)

type TaskController struct {
	beego.Controller
}

// @Title set task
// @Description 列表上課內容
// @
// @Success  200  object  models.IClass  "上課資訊"
// @
// @Resource 關於上課內容
// @Router / [get]
func (c *TaskController) Get() {
	// check path
	var result [5]models.TaskUrl
	filename := "./data/task.json"
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0777)
	}

	if _, err := os.Stat(filename); err == nil {
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

// @Title set task
// @Description 列表上課內容
// @
// @Success  200  object  models.IClass  "上課資訊"
// @
// @Resource 關於上課內容
// @Router / [post]
func (c *TaskController) Post() {
	var taskLists []models.TaskUrl
	json.Unmarshal([]byte(c.GetString("taskList")), &taskLists)

	// 檢查資料夾
	filename := "./data/task.json"
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0777)
	}

	// write json
	content, err := json.MarshalIndent(taskLists, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	_ = os.WriteFile(filename, content, 0777)

	c.Data["json"] = &taskLists
	c.ServeJSON()
}
