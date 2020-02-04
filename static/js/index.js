window.onload = function() {
    Vue.component('crawler-template', {
        props: ['item'],
        template: '#crawler-template'
    })

    var app = new Vue({
        el: '#app',
        data: {
            pages: 1,
            isCrawlerTable: false,
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
                var self = this;
                $.ajax({
                    type: 'POST',
                    url: '/',
                    data: {},
                    success: function(data) {
                        self.classes = data;
                        self.isCrawlerTable = true;
                    }
                });
            }
        }
    })
}