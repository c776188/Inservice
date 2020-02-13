window.onload = function() {
    new Vue({
        el: "#app",
        vuetify: new Vuetify(),
        data: {
            loading: true,
            search: "",
            headers: [{
                    text: "ID",
                    align: "left",
                    sortable: false,
                    value: "id"
                },
                { text: "課程名稱", value: "Name" },
                { text: "開課地點", value: "Detail.Location" },
                { text: "距離時間", value: "Detail.MapElement.duration.text" },
                { text: "報名時間", value: "Detail.SignUpStatus" },
                { text: "上課日期", value: "Detail.AttendClassTime" },
                { text: "研習時數", value: "Detail.StudyHours" },
                { text: "登錄日期", value: "Detail.EntryDate" }
            ],
            classes: [],
            urlPrefix: "https://www1.inservice.edu.tw/NAPP/CourseView.aspx?cid="
        },
        created: function() {
            this.callCrawler();
        },
        methods: {
            callCrawler() {
                this.loading = true;
                let self = this;

                axios.post("/")
                    .then(res => {
                        self.classes = res.data;
                        self.loading = false;
                    })
                    .catch(error => {
                        console.error(error);
                    });
            }
        }
    });
};