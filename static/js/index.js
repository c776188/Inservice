window.onload = function() {
    new Vue({
        el: '#app',
        vuetify: new Vuetify(),
        data: {
            isloading: true,
            search: '',
            headers: [{
                    text: 'ID',
                    align: 'left',
                    sortable: false,
                    value: 'ID',
                },
                { text: '課程名稱', value: 'Name' },
                { text: '開課地點', value: 'Detail.Location' },
                { text: '距離時間', value: 'Detail.MapElement.duration.text' },
                { text: '報名時間', value: 'Detail.SignUpStatus' },
                { text: '上課日期', value: 'Detail.AttendClassTime' },
                { text: '研習時數', value: 'Detail.StudyHours' },
                { text: '登錄日期', value: 'Detail.EntryDate' },
            ],
            classes: [
                {}
            ]
        },
        created: function() {
            this.callCrawler()
        },
        methods: {
            callCrawler() {
                this.isloading = true;
                var self = this;
                $.ajax({
                    type: 'POST',
                    url: '/',
                    data: {},
                    success: function(data) {
                        self.classes = data;
                        self.isloading = false;
                    }
                });
            }
        }
    })
}