<template>
    <div id="app">
        <el-menu theme="dark" class="el-menu-demo" mode="horizontal">
            <el-menu-item index="1"><span style="font-family: 'Satisfy', cursive; font-size: 30px;">Less</span>
            </el-menu-item>
            <a href="https://github.com/lnquy/less" class="menu-right" target="_blank">Github</a>
            <div class="centered-menu-block">
                <el-date-picker v-model="date_picker" type="date" placeholder="Pick a date"></el-date-picker>
                <el-button type="primary" @click="getTrendingRepos()" :loading="loading">Go!</el-button>
            </div>
        </el-menu>
        <el-row>
            <el-col :span="16" class="repos-col">
                <div v-if="repos.length != 0" >
                    <h2>Github Trending</h2>
                    <div class="repos">
                        <repository v-for="repo in repos" :repo="repo" :key="repo.name"></repository>
                    </div>
                </div>
                <div v-else style="padding: 30px">
                    No data
                    <div><br>
                        <h4>Serverless website on Amazon Web Services</h4>
                        <img src="./assets/less-arch.jpg" style="width: 70%">
                    </div>
                </div>
            </el-col>
        </el-row>

        <div class="footer">
            Built in love with Go. <a href="https://github.com/lnquy" target="_blank" style="text-decoration: none">@lnquy</a>
        </div>
    </div>
</template>

<script>
    import {global} from "./main.js";
    import Repository from "./components/Repository.vue";

    export default {
        data() {
            return {
                loading: false,
                date_picker: new Date().toISOString().slice(0, 10),
                repos: [],
            }
        },
        components: {
            "repository": Repository,
        },
        methods: {
            getTrendingRepos() {
                this.loading = true;
                this.$http.post(global.getCatererUrl(), "{\"date\": \"" + this.getUTCDate() + "\"}").then(resp => {
                    if (resp.body != null && resp.body !== "") {
                        this.repos = JSON.parse(resp.body);
                    } else {
                        this.repos = [];
                    }
                    this.loading = false;
                }, err => {
                    this.repos = [];
                    this.loading = false;
                    this.$notify.error({
                        title: 'Error',
                        message: "Failed to get trending repositories!",
                        duration: 5000
                    });
                });
            },
            getUTCDate() {
                let nd = new Date(this.date_picker);
                return nd.getFullYear() + "-" + ("0" + (nd.getMonth() + 1)).slice(-2) + "-" + ("0" + nd.getDate()).slice(-2)
            },
        },
    }
</script>

<style>
    body {
        margin: 0;
        font-family: Helvetica, sans-serif;
    }

    span, div {
        margin: 0;
        padding: 0;
    }

    .footer {
        bottom: 0;
        font-size: 14px;
        color: #BBB;
        text-align: center;
        padding: 20px;
    }

    .el-menu,
    .el-menu .el-menu-item {
        border-radius: 0;
        z-index: 3;
    }
</style>

<style scoped>
    .centered-menu-block{
        position: absolute;
        padding-top: 11px;
        width: 100%;
        text-align: center;
        z-index: 1;
    }

    .menu-right {
        position: absolute;
        right: 20px;
        top: 20px;
        color: #bfcbd9;
        text-decoration: none;
        z-index: 2;
    }
    .menu-right:hover {
        color: #009DFB;
    }

    .repos-col {
        width: 100%;
        text-align: center;
    }

    .repos {
        max-width: 600px;
        text-align: center;
        margin: 0 auto;
    }
</style>
