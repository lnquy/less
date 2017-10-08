import Vue from 'vue'
import VueResource from 'vue-resource'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-default/index.css'
import locale from 'element-ui/lib/locale/lang/en'
import App from './App.vue'

Vue.use(VueResource);
Vue.use(ElementUI, { locale });

Vue.config.productionTip = false;
Vue.http.options.xhr = {withCredentials: true};

new Vue({
    el: '#app',
    render: h => h(App)
});

export const global = {
    getCatererUrl() {
        return "https://6epko5iya8.execute-api.ap-southeast-1.amazonaws.com/dev/api/v1/trending";
    }
};
