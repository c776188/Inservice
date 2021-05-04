package controllers

import (
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
// @Router / [post]
func (c *TaskController) Post() {
	var result []string

	c.Data["json"] = &result
	c.ServeJSON()
}
