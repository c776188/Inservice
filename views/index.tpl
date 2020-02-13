<!DOCTYPE html>
<html>

<head>
    <title>Crawler</title>
    <link href="https://fonts.googleapis.com/css?family=Roboto:100,300,400,500,700,900" rel="stylesheet" />
    <link href="https://cdn.jsdelivr.net/npm/@mdi/font@4.x/css/materialdesignicons.min.css" rel="stylesheet" />
    <link href="https://cdn.jsdelivr.net/npm/vuetify@2.x/dist/vuetify.min.css" rel="stylesheet" />
    <link href="https://use.fontawesome.com/releases/v5.0.13/css/all.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/font-awesome@4.x/css/font-awesome.min.css" rel="stylesheet">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Material+Icons|Material+Icons+Outlined|Material+Icons+Two+Tone|Material+Icons+Round|Material+Icons+Sharp" />
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no, minimal-ui" />

    <style type="text/css">
        [v-cloak] {
            display: none;
        }
    </style>
</head>

<body>
    <div id="app">
        <v-app id="inspire">
            <!-- <div v-if="loading" id="loading">Loading~~~~</div> -->
            <v-card dark>
                <v-card-title>
                    全國教師在職進修網
                    <v-spacer></v-spacer>
                    <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details>
                    </v-text-field>
                </v-card-title>
                <v-data-table :headers="headers" :items="classes" :search="search" dark :loading="loading" loading-text="Loading... Please wait">
                    <template slot="item.ID" slot-scope="props">
                        <a :href="urlPrefix + props.value" target="_blank" :title="props.value">{{props.value}}</a>
                    </template>
                </v-data-table>
            </v-card>
        </v-app>
    </div>

    <script src="https://code.jquery.com/jquery-3.4.1.js" integrity="sha256-WpOohJOqMqqyKL9FccASB9O0KwACQJpFTUBLTYOVvVU=" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/vuetify@2.x/dist/vuetify.js"></script>
    <script src="/static/js/index.js"></script>
</body>

</html>