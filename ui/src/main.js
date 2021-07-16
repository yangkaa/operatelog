/*
 * @Author: your name
import ElementUI from "element-ui"; // 2.1引入结构
import "element-ui/lib/theme-chalk/index.css"; // 2.2引入样式
import Vue from "vue";
import App from "./App";
import router from "./router";
 */
// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.

import ElementUI from "element-ui"; // 2.1引入结构
import "element-ui/lib/theme-chalk/index.css"; // 2.2引入样式
import Moment from "moment";
import Vue from "vue";
import App from "./App";
import router from "./router";
Vue.use(ElementUI);
Vue.config.productionTip = false;
Vue.prototype.moment = Moment;

/* eslint-disable no-new */
new Vue({
  el: "#app",
  router,
  components: { App },
  template: "<App/>"
});
