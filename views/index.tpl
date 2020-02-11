<!DOCTYPE html>
<html>

<head>
    <title>Crawler</title>
    <script src="https://code.jquery.com/jquery-3.4.1.js" integrity="sha256-WpOohJOqMqqyKL9FccASB9O0KwACQJpFTUBLTYOVvVU=" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>
    <script src="/static/js/index.js"></script>

    <style type="text/css">
        [v-cloak] {
            display: none;
        }
    </style>
</head>

<body>
    <div id="app">
        <div v-if="isloading" id="loading"> Loading~~~~</div>
        <table v-cloak v-if="isCrawlerTable" cellspacing="0" cellpadding="4" border="1">
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Location</th>
                    <th>Duration</th>
                    <th>SignUpTime</th>
                    <th>AttendClassTime</th>
                    <th>StudyHours</th>
                </tr>
            </thead>
            <tbody id="rows">
                <tr is="crawler-template" v-for="(item, index) in classes" :item="item" :key="index"></tr>
            </tbody>
        </table>
    </div>
</body>

<script type="text/x-template" id="crawler-template">
    <tr>
        <td>{{item.ID}}</td>
        <td>{{item.Name}}</td>
        <!-- <td>{{item.Location}}</td> -->
        <td>{{item.Detail.Location}}</td>
        <td>
            <span v-for="r in item.Detail.MapDetail.rows">
                <span v-for="e in r.elements">
                    {{e.duration.text}}
                </span>
            </span>
        </td>
        <td>{{item.Detail.SignUpStatus}} <br/> {{item.Detail.SignUpTime}}</td>
        <td>{{item.Detail.AttendClassTime}}</td>
        <td>{{item.Detail.StudyHours}}</td>
    </tr>
</script>

</html>