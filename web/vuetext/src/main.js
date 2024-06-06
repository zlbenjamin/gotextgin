import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconVue from '@element-plus/icons-vue'


const app = createApp(App)

app.use(router)

app.use(ElementPlus)

// 全局注册
for (const [key, component] of Object.entries(ElementPlusIconVue)) {
    app.component(key, component)
}

app.mount('#app')
