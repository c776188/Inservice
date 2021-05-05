window.onload = function () {
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
                    value: "ID"
                },
                {
                    text: "課程名稱",
                    value: "Name"
                },
                {
                    text: "開課地點",
                    value: "Detail.Location"
                },
                {
                    text: "距離時間",
                    value: "Detail.MapElement.duration.text"
                },
                {
                    text: "報名時間",
                    value: "Detail.SignUpTime"
                },
                {
                    text: "上課日期",
                    value: "Detail.AttendClassTime"
                },
                {
                    text: "研習時數",
                    value: "Detail.StudyHours"
                },
                {
                    text: "登錄日期",
                    value: "Detail.EntryDate"
                },
                {
                    text: "連結",
                    sortable: false,
                    value: "url"
                }
            ],
            classes: [],
            urlPrefix: "https://www1.inservice.edu.tw/NAPP/CourseView.aspx?cid=",
            selectedHeaders: [],
            showHeaders: [],
            dialog: false,
            taskList: [],
        },
        created: function () {
            this.selectedHeaders = this.headers.slice()
            this.showHeaders = this.headers.slice()
            this.callCrawler();
        },
        computed: {
            likesAllFruit() {
                return this.selectedHeaders.length === this.headers.length
            },
            likesSomeFruit() {
                return this.selectedHeaders.length > 0 && !this.likesAllFruit
            },
            icon() {
                if (this.likesAllFruit) return 'mdi-close-box'
                if (this.likesSomeFruit) return 'mdi-minus-box'
                return 'mdi-checkbox-blank-outline'
            },
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

                axios.get("/task")
                    .then(res => {
                        self.taskList = res.data;
                    })
                    .catch(error => {
                        console.error(error);
                    });
            },
            setSelected() {
                this.showHeaders = [];
                for (let i = 0; i < this.headers.length; i++) {
                    if (selectedHeaders.indexOf(this.headers[i].value)) {
                        this.showHeaders.push(this.headers[j]);
                        break;
                    }
                }
            },
            toggle() {
                this.$nextTick(() => {
                    if (this.likesAllFruit) {
                        this.selectedHeaders = []
                        this.showHeaders = []
                    } else {
                        this.selectedHeaders = this.headers.slice()
                        this.showHeaders = this.headers.slice()
                    }
                })
            },
            gotoLink(url) {
                window.open(url, '_blank');
            },
            storeTask() {
                let self = this;

                // 課程名字
                var count = 1;
                this.classes.forEach(function (item) {
                    for (var i = 0; i < 5; i++) {
                        if (self.taskList[i].ID == item.ID) {
                            self.taskList[i].Name = item.Name;
                            count++;
                        }

                        if (count == 5) {
                            break;
                        }
                    }
                });

                var bodyFormData = new FormData();
                bodyFormData.append('taskList', JSON.stringify(this.taskList));

                axios({
                        method: "post",
                        url: "/task",
                        data: bodyFormData,
                        headers: {
                            "Content-Type": "multipart/form-data"
                        },
                    })
                    .then(res => {
                        console.log(res.data)
                    })
                    .catch(error => {
                        console.error(error);
                    });

                this.dialog = false
            }
        }
    });
};