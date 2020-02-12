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
                { text: 'Name', value: 'Name' },
                { text: 'Location', value: 'Detail.Location' },
                { text: 'Duration', value: 'Detail.MapDetail.rows[0].elements[0].duration.text' },
                { text: 'SignUpTime', value: 'Detail.SignUpStatus' },
                { text: 'AttendClassTime', value: 'Detail.AttendClassTime' },
                { text: 'StudyHours', value: 'Detail.StudyHours' },
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