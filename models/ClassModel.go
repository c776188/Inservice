package models

// 課程資訊
type IClass struct {
	ID       string
	Name     string
	Location string
	Detail   IDetail
}

// 課程詳細資訊
type IDetail struct {
	SignUpStatus    string
	SignUpTime      string
	AttendClassTime string
	StudyHours      string
	Location        string
	EntryDate       string
	MapElement      Elements
}

// google map 串接資訊
type GMap struct {
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
