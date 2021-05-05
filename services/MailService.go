package services

import (
	"Inservice/models"
	"fmt"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/utils"
)

func Mail(task models.TaskUrl) {
	mapConfig, _ := config.NewConfig("ini", "conf/env.conf")
	// searchUrl := mapConfig.String("SEARCH_URL")
	config := `{"username":"` + mapConfig.String("SERVER_MAIL") + `","password":"` + mapConfig.String("SERVER_MAIL_PASSWORD") + `","host":"smtp.gmail.com","port":587}`
	email := utils.NewEMail(config)
	email.To = []string{mapConfig.String("MAIL_TARGET")}
	email.From = mapConfig.String("SERVER_MAIL")
	email.Subject = "Inservice Alert"
	email.HTML = "<h1>" + task.Name + " 已開放</h1><br/> <a href='https://www1.inservice.edu.tw/NAPP/CourseView.aspx?cid=" + task.ID + "'>課程連結</a>"
	err := email.Send()
	if err != nil {
		fmt.Println(err)
		return
	}
}
