window.onload = function() {
    Vue.component('crawler-template', {
        props: ['item'],
        template: '#crawler-template'
    })

    var app = new Vue({
        el: '#app',
        data: {
            isCrawlerTable: false,
            isloading: true,
            classes: [
                {}
            ]
        },
        created: function() {
            this.callCrawler()
        },
        methods: {
            callCrawler() {
                this.isCrawlerTable = false;
                this.isloading = true;
                var self = this;
                $.ajax({
                    type: 'POST',
                    url: '/',
                    data: {},
                    success: function(data) {
                        self.classes = data;
                        self.isCrawlerTable = true;
                        self.isloading = false;
                    }
                });
            }
        }
    })
}