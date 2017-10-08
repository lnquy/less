<template>
    <el-card class="box-card">
        <div slot="header" class="clearfix">
            <span><a :href="repo.url | github" target="_blank">{{ repo.name }}</a></span>
            <a :href="repo.url | github" style="float: right;" target="_blank"><i class="el-icon-star-off" style="color: #CCC"></i></a>
        </div>
        <div class="card-content">
            <p class="description">{{ repo.description | checkNull }}</p>
            <div style="font-size: 13px; color: #CCC">
                <span>{{ repo.language | checkNull }}</span><span v-if="repo.language != 'null'">&nbsp;&nbsp;&nbsp;</span>
                <span><i class="el-icon-star-on"></i> {{ repo.stars }}</span>&nbsp;&nbsp;&nbsp;
                <span><i class="el-icon-share"></i> {{ repo.forks }}</span>
                <span style="float: right">{{ repo.today_stars }} stars today</span>
            </div>
        </div>
    </el-card>
</template>

<script>
    export default {
        data() {
            return {}
        },
        props: ["repo"],
        filters: {
            github(val) {
                return "https://github.com" + val;
            },
            checkNull(val) {
                return (val == "null") ? "" : val;
            }
        }
    }
</script>

<style>
    p.description {
        margin-top: 0;
        font-size: 16px;
        line-height: 2em;
    }
    .clearfix:before,
    .clearfix:after {
        display: table;
        content: "";
    }
    .clearfix:after {
        clear: both
    }
    .clearfix > span > a {
        line-height: 24px;
        font-size: 20px;
        font-weight: 600;
        text-decoration: none;
    }

    .box-card {
        text-align: left;
        max-width: 600px;
        margin: 10px;
    }
</style>
