import Vue from "vue"
import "@/styles/index.scss"
// import "@/assets/iconfont/iconfont"
// import '@/assets/iconfont/iconfont.css'
import App from "./App.vue"
import Fit2CloudUI from "fit2cloud-ui"
import ElementUI from "element-ui"
import i18n from "./i18n"
import router from "./router"
import store from './store'
// import icons from './icons'
// import plugins from "./plugins";
// import directives from "./directive";
// import filters from "./filters";
import "./permission"

Vue.config.productionTip = false

Vue.use(ElementUI, {
  size: "small",
  i18n: (key, value) => i18n.t(key, value)
})
Vue.use(Fit2CloudUI, {
  i18n: (key, value) => i18n.t(key, value)
})


new Vue({
  el: '#app',
  i18n,
  router,
  store,
  render: h => h(App),
})
