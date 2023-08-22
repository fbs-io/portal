import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import base from './base'
import i18n from './locales'
import store from './stores'

const app = createApp(App)

app.use(createPinia())
app.use(store)
app.use(router)
app.use(ElementPlus)
app.use(base)
app.use(i18n)
app.mount('#app')
