// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import iView from 'iview';
import 'iview/dist/styles/iview.css';
import './assets/css/base.css'

import VueI18n from 'vue-i18n';
import en from 'iview/dist/locale/en-US';
import zh from 'iview/dist/locale/zh-CN';

import myen from './i18n/en-US';
import myzh from './i18n/zh-CN';


Vue.use(VueI18n);

Vue.config.productionTip = false

Vue.use(iView, {
  transfer: true
});


// 语言
Vue.locale = () => {};

const messages = {
    en: Object.assign(myen, en),
    zh: Object.assign(myzh, zh)
};

let mylang = localStorage.getItem("etcd-language") || 'en';

// Create VueI18n instance with options
const i18n = new VueI18n({
    locale: mylang,  // set locale
    messages  // set locale messages
});

// 判断字符串是否是以xx为前缀
Vue.prototype.HasPrefix = (str, prefix) => {
  console.log(str, prefix,str.substring(str.length - prefix.length));
  
  if(prefix == null || prefix == "" || str.length == 0 || prefix.length > str.length)
     return false
  if(str.substring(0, prefix.length) == prefix)
     return true
  else
     return false
}


/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  i18n,
  components: {
    App
  },
  template: '<App/>'
})
