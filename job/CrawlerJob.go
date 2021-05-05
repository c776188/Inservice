package job

import (
	"Inservice/models"
	"Inservice/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/toolbox"
)

func CrawlerTaskListener() {
	tk := toolbox.NewTask("crawlerTask", "0 0/5 * * * *", func() error { readTaskFile(); return nil })
	err := tk.Run()
	if err != nil {
		fmt.Println(err)
	}
	toolbox.AddTask("crawlerTask", tk)
	toolbox.StartTask()
}

func readTaskFile() {
	// check path
	var contents []models.TaskUrl
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

		_ = json.Unmarshal([]byte(file), &contents)

		for i, single := range contents {
			checkStatus(single)

			// 超過五篇則無效
			if i >= 5 {
				break
			}
		}
	}
}

func checkStatus(task models.TaskUrl) {
	url := "https://www1.inservice.edu.tw/NAPP/CourseView.aspx?cid=" + task.ID
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "www1.inservice.edu.tw")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// goquery 爬蟲取得資訊
	// var status string
	// 報名狀態
	dom.Find("#ctl00_CPH_Content_RLB_OnlineApply_w").Each(func(i int, selection *goquery.Selection) {
		// status = selection.Text()
		class, _ := selection.Attr("class")
		// fmt.Println(class, selection.Text())
		if strings.Contains(class, "rbDisabled") {
			fmt.Println("還沒開放")
		} else {
			fmt.Println("已開放")
			services.Mail(task)
		}
	})
}
