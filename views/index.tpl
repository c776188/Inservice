<!DOCTYPE html>
<html>

<head>
    <title>Inservice</title>
    <link href="https://fonts.googleapis.com/css?family=Roboto:100,300,400,500,700,900" rel="stylesheet" />
    <link href="https://cdn.jsdelivr.net/npm/@mdi/font@4.x/css/materialdesignicons.min.css" rel="stylesheet" />
    <link href="https://cdn.jsdelivr.net/npm/vuetify@2.x/dist/vuetify.min.css" rel="stylesheet" />
    <link href="https://use.fontawesome.com/releases/v5.0.13/css/all.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/font-awesome@4.x/css/font-awesome.min.css" rel="stylesheet">
    <link rel="stylesheet"
        href="https://fonts.googleapis.com/css?family=Material+Icons|Material+Icons+Outlined|Material+Icons+Two+Tone|Material+Icons+Round|Material+Icons+Sharp" />
    <meta name="viewport"
        content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no, minimal-ui" />

    <style type="text/css">
        [v-cloak] {
            display: none;
        }

        .v-list-item--link::before {
            background-color: red;
        }
    </style>
</head>

<body>
    <div id="app">
        <v-app id="inspire">
            <!-- <div v-if="loading" id="loading">Loading~~~~</div> -->
            <v-card dark height="100%">
                <v-card-title>
                    全國教師在職進修網
                    <v-spacer></v-spacer>
                    <v-select v-model="selectedHeaders" item-value="value" :items="headers" label="Favorite Header"
                        @input="setSelected" multiple dark>
                        <template v-slot:prepend-item>
                            <v-list-item ripple @click="toggle">
                                <v-list-item-action>
                                    <v-icon :color="selectedHeaders.length > 0 ? 'indigo darken-4' : ''">{{ icon }}
                                    </v-icon>
                                </v-list-item-action>
                                <v-list-item-content>
                                    <v-list-item-title>Select All</v-list-item-title>
                                </v-list-item-content>
                            </v-list-item>
                            <v-divider class="mt-2"></v-divider>
                        </template>
                    </v-select>
                    <v-spacer></v-spacer>
                    <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details>
                    </v-text-field>
                    <v-spacer></v-spacer>
                    <!-- <v-btn class="ma-2" color="blue" @click="setTask" dark>
                        <v-icon dark>
                            mdi-wrench
                        </v-icon>
                    </v-btn> -->
                    <template>
                        <v-row justify="center">
                            <v-dialog v-model="dialog" persistent max-width="600px">
                                <template v-slot:activator="{ on, attrs }">
                                    <v-btn color="primary" dark v-bind="attrs" v-on="on">
                                        Task
                                    </v-btn>
                                </template>
                                <v-card dark>
                                    <v-card-title>
                                        <span class="headline">Task List</span>
                                    </v-card-title>
                                    <v-card-text>
                                        <v-container>
                                            <v-row v-for="n in 5">
                                                <v-col cols="12" sm="3" md="3">
                                                    <!-- v-model='taskList[$n-1].ID' -->
                                                    <v-text-field label="ID" hint="請輸入ID" v-model="taskList[n-1].ID"></v-text-field>
                                                </v-col>
                                                <v-col cols="12" sm="9" md="9">
                                                    <span style="font-size: 20px;">課程名稱: {{taskList[n-1].Name}}</span>
                                                </v-col>
                                            </v-row>
                                        </v-container>
                                    </v-card-text>
                                    <v-card-actions>
                                        <v-spacer></v-spacer>
                                        <v-btn color="blue darken-1" text @click="dialog = false">
                                            Close
                                        </v-btn>
                                        <v-btn color="blue darken-1" text @click="storeTask">
                                            Save
                                        </v-btn>
                                    </v-card-actions>
                                </v-card>
                            </v-dialog>
                        </v-row>
                    </template>
                </v-card-title>
                <v-data-table :headers="showHeaders" :items="classes " :search="search " :loading="loading "
                    loading-text="Loading... Please wait ">
                    <template v-slot:item.url="{ item } ">
                        <v-icon small @click="gotoLink(urlPrefix + item.ID)">fas fa-external-link-alt</v-icon>
                    </template>
                </v-data-table>
            </v-card>
        </v-app>
    </div>

    <script src="https://unpkg.com/axios/dist/axios.min.js "></script>
    <script src="https://code.jquery.com/jquery-3.4.1.js "
        integrity="sha256-WpOohJOqMqqyKL9FccASB9O0KwACQJpFTUBLTYOVvVU=" crossorigin=" anonymous "></script>
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js "></script>
    <script src="https://cdn.jsdelivr.net/npm/vuetify@2.x/dist/vuetify.js "></script>
    <script src="/static/js/index.js "></script>
</body>

</html>