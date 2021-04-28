package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	_ "github.com/astaxie/beego/config"
)

type defaultKey struct {
	EVENTVALIDATION string
	VIEWSTATE       string
	Class           []iClass
}

// 課程資訊
type iClass struct {
	ID       string
	Name     string
	Location string
	Detail   iDetail
}

// 課程詳細資訊
type iDetail struct {
	SignUpStatus    string
	SignUpTime      string
	AttendClassTime string
	StudyHours      string
	Location        string
	EntryDate       string
	MapElement      Elements
}

// google map 串接資訊
type gMap struct {
	Destination_addresses []string `json:"destination_addresses"`
	Origin_addresses      []string `json:"origin_addresses"`
	Rows                  []Row    `json:"rows"`
	Status                string   `json:"status"`
}

type Row struct {
	Elements []Elements `json:"elements"`
}

type Elements struct {
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
	Fare     Fare     `json:"fare"`
	Status   string   `json:"status"`
}

type Distance struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Duration struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Fare struct {
	Currency string `json:"currency"`
	Text     string `json:"text"`
	Value    int    `json:"value"`
}

var searchUrl = ""
var cookie = ""

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
// @Success  200  object  iClass  "上課資訊"
// @
// @Resource 關於上課內容
// @Router / [post]
func (c *MainController) Post() {
	searchUrl = "https://www1.inservice.edu.tw/script/IndexQuery.aspx?city=9"

	var result []iClass

	// 爬蟲
	result = getInitInservice().Class

	// 取得map距離
	// result = getMapDuration(result)

	c.Data["json"] = &result
	c.ServeJSON()
}

func getMapDuration(result []iClass) []iClass {
	// 計算距離
	locations := []string{}
	var locationString string
	maxMapCount := 25
	nloop := 0
	for i, r := range result {
		// push array
		locations = append(locations, r.Detail.Location)

		if i%maxMapCount == maxMapCount-1 || i == len(result)-1 {
			// handle string and call map api
			locationString = TrimSpaceNewlineInString(strings.Join(locations[:], "|"))
			tmpMap := callMapDistanceMatrix(locationString)

			// assign element
			for j, element := range tmpMap.Rows[0].Elements {
				result[j+nloop*maxMapCount].Detail.MapElement = element
			}

			// reset data
			locations = nil
			nloop++
		}
	}

	return result
}

func TrimSpaceNewlineInString(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, "")
}

// 取得key
func getInitInservice() defaultKey {
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// catch cookie
	for _, c := range res.Cookies() {
		if c.Name == "ASP.NET_SessionId" {
			cookie = c.Value
		}
	}

	// goquery 爬蟲取得資訊
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result defaultKey

	// 更新key
	result = updateKey(dom, result)

	// 搜尋資訊加密
	searchKey := make(map[string]string)
	searchKey["__VIEWSTATEGENERATOR"] = "82F443D6"
	searchKey["__VIEWSTATEENCRYPTED"] = ""
	searchKey["__EVENTARGUMENT"] = ""
	searchKey["__LASTFOCUS"] = ""
	searchKey["ddlQueryType"] = "byCity"
	searchKey["ddlCityList"] = "9"
	searchKey["ddlSchoolLevelByCity"] = "50"
	searchKey["ddlCourseTag"] = ""

	// 取得資訊
	result.Class = postSearchInservice(result, encodeSendData(searchKey))

	return result
}

// 搜尋資料
func postSearchInservice(key defaultKey, searchKey string) []iClass {
	postData := make(map[string]string)
	postData["__EVENTVALIDATION"] = key.EVENTVALIDATION
	postData["__VIEWSTATE"] = key.VIEWSTATE
	postData["__EVENTTARGET"] = "dgSelectResult$_ctl24$_ctl0"
	postData["Button1"] = "查詢"
	payload := strings.NewReader(searchKey + encodeSendData(postData))

	req, _ := http.NewRequest("POST", searchUrl, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "ASP.NET_SessionId="+cookie)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "www1.inservice.edu.tw")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "19788")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	// goquery 爬蟲取得資訊
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// 更新key
	key = updateKey(dom, key)

	return postInservicePage(0, key, searchKey)
}

// 爬蟲 取得課程資訊
func postInservicePage(page int, key defaultKey, searchKey string) []iClass {
	postData := make(map[string]string)
	postData["__EVENTVALIDATION"] = key.EVENTVALIDATION
	postData["__VIEWSTATE"] = key.VIEWSTATE
	postData["__EVENTTARGET"] = "dgSelectResult$_ctl24$_ctl" + strconv.Itoa(page)
	payload := strings.NewReader(searchKey + encodeSendData(postData))

	req, _ := http.NewRequest("POST", searchUrl, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "ASP.NET_SessionId="+cookie)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "www1.inservice.edu.tw")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "19788")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	// goquery 爬蟲取得資訊
	classes := []iClass{}
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// 更新key
	key = updateKey(dom, key)

	// 取得資訊
	dom.Find(".cinfo-r1>tbody>tr>td:first-child").Each(func(i int, selection *goquery.Selection) {
		temp := iClass{ID: selection.Text(), Detail: postInserviceDetail(selection.Text())}
		classes = append(classes, temp)
	})

	dom.Find(".cinfo-r2").Each(func(i int, selection *goquery.Selection) {
		classes[i].Name = selection.Text()
	})

	dom.Find(".cinfo-r3").Each(func(i int, selection *goquery.Selection) {
		classes[i].Location = selection.Text()
	})

	nextPage := false
	dom.Find(".cssctsTitle2>td>a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if !ok {
			fmt.Println("error")
		}

		if regexp.MustCompile("ctl" + strconv.Itoa(page+1) + "'").Match([]byte(href)) {
			nextPage = true
		}
	})

	if nextPage {
		return append(classes, postInservicePage(page+1, key, searchKey)...)
	}

	return classes
	// fmt.Println(res)
}

// 更新取得key
func updateKey(dom *goquery.Document, key defaultKey) defaultKey {
	dom.Find("input#__VIEWSTATE").Each(func(i int, selection *goquery.Selection) {
		key.VIEWSTATE, _ = selection.Attr("value")
	})

	dom.Find("input#__EVENTVALIDATION").Each(func(i int, selection *goquery.Selection) {
		key.EVENTVALIDATION, _ = selection.Attr("value")
	})

	return key
}

// encode map to string
func encodeSendData(m map[string]string) string {
	b := new(bytes.Buffer)
	for k, v := range m {
		fmt.Fprintf(b, "%s=%s&", k, url.QueryEscape(v))
	}
	return b.String()
}

// 取得詳細資料
func postInserviceDetail(id string) iDetail {
	url := "https://www1.inservice.edu.tw/NAPP/CourseView.aspx?cid=" + id

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
	var detail iDetail
	// 報名狀態
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(7) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c4.cs1").Each(func(i int, selection *goquery.Selection) {
		detail.SignUpStatus = TrimSpaceNewlineInString(selection.Text())
	})

	// 報名時間
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(7) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c2.cs1").Each(func(i int, selection *goquery.Selection) {
		// 判斷狀態是否顯示更詳細報名時間
		layout := "2006/01/0215:04"
		restime := strings.Split(detail.SignUpStatus, "開放線上報名")
		restime = strings.Split(restime[1], "起")
		t, err := time.Parse(layout, restime[0])
		if err != nil {
			detail.SignUpTime = TrimSpaceNewlineInString(selection.Text())
		} else {
			detail.SignUpTime = t.Format("2006/01/02 15:04")
		}
	})

	// 上課日期
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c2.cs1").Each(func(i int, selection *goquery.Selection) {
		tempAttendClassTime := TrimSpaceNewlineInString(selection.Text())

		// 日期同一天的話，去掉後面
		splitAttend := strings.Split(tempAttendClassTime, "至")
		if splitAttend[0] == splitAttend[1] {
			detail.AttendClassTime = splitAttend[0]
		} else {
			detail.AttendClassTime = tempAttendClassTime
		}
	})

	// 研習時數
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(3) > td.cw_table_info.c2.cs3").Each(func(i int, selection *goquery.Selection) {
		tmpStudyHours := TrimSpaceNewlineInString(selection.Text())
		tmpStudyHours = strings.ReplaceAll(tmpStudyHours, "/學分", "")
		tmpStudyHours = strings.ReplaceAll(tmpStudyHours, "小時", "H")
		detail.StudyHours = tmpStudyHours
	})

	// 開課地點
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(2) > td.cw_table_info.c2.cs3").Each(func(i int, selection *goquery.Selection) {
		detail.Location = TrimSpaceNewlineInString(selection.Text())
	})

	// 登錄日期
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(3) > div.cid_block > div:nth-child(2)").Each(func(i int, selection *goquery.Selection) {
		detail.EntryDate = strings.ReplaceAll(selection.Text(), "．登錄日期：", "")
	})
	return detail
}

// 取得map資料
func callMapDistanceMatrix(destinations string) gMap {
	// 從config取得map key
	mapConfig, err := config.NewConfig("ini", "conf/env.conf")
	mapKey := mapConfig.String("gMapKey")

	sendData := make(map[string]string)
	sendData["units"] = "imperial"
	sendData["origins"] = mapConfig.String("origins")
	sendData["destinations"] = destinations
	sendData["mode"] = "driving"
	sendData["key"] = mapKey
	url := "https://maps.googleapis.com/maps/api/distancematrix/json?" + encodeSendData(sendData)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// json string to struct
	g := gMap{}
	err2 := json.Unmarshal(body, &g)
	if err2 != nil {
		log.Fatalln(err2)
	}

	return g
}
