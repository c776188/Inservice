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
	MapDetail       gMap
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

var searchUrl = "https://www1.inservice.edu.tw/script/IndexQuery.aspx?city="

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

// post 取得課程資訊
func (c *MainController) Post() {
	searchUrl = searchUrl + "9"

	var result []iClass

	// post
	result = getInitInservice().Class

	c.Data["json"] = &result
	c.ServeJSON()
}

// 取得key
func getInitInservice() defaultKey {
	req, _ := http.NewRequest("GET", searchUrl, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "www1.inservice.edu.tw")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Cookie", "ASP.NET_SessionId=lhzlilg1z2e0ibwneqi1keex")
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

	var result defaultKey

	// 更新key
	result = updateKey(dom, result)

	// 搜尋資訊加密
	searchKey := make(map[string]string)
	searchKey["__VIEWSTATEGENERATOR"] = "82F443D6"
	searchKey["__VIEWSTATEENCRYPTED"] = ""
	searchKey["__EVENTARGUMENT"] = ""
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
	postData["button1"] = "查詢"
	payload := strings.NewReader(searchKey + encodeSendData(postData))

	req, _ := http.NewRequest("POST", searchUrl, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "ASP.NET_SessionId=lhzlilg1z2e0ibwneqi1keex")
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
	req.Header.Add("Cookie", "ASP.NET_SessionId=lhzlilg1z2e0ibwneqi1keex")
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
		detail.SignUpStatus = selection.Text()
	})

	// 報名時間
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(7) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c2.cs1").Each(func(i int, selection *goquery.Selection) {
		detail.SignUpTime = selection.Text()
	})

	// 上課日期
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c2.cs1").Each(func(i int, selection *goquery.Selection) {
		detail.AttendClassTime = selection.Text()
	})

	// 研習時數
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(3) > td.cw_table_info.c2.cs3").Each(func(i int, selection *goquery.Selection) {
		detail.StudyHours = selection.Text()
	})

	// 開課地點
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(2) > td.cw_table_info.c2.cs3").Each(func(i int, selection *goquery.Selection) {
		detail.Location = selection.Text()
		// detail.MapDetail = getMapDuration(selection.Text())
	})

	return detail
}

// 取得map資料
func getMapDuration(destinations string) gMap {
	// 從config取得map key
	mapConfig, err := config.NewConfig("ini", "conf/env.conf")
	mapKey := mapConfig.String("gMapKey")

	sendData := make(map[string]string)
	sendData["units"] = "imperial"
	sendData["origins"] = mapConfig.String("origins")
	sendData["destinations"] = destinations
	sendData["mode"] = "transit"
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
