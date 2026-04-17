import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter } from './router'
import TDesign from 'tdesign-vue-next'
import 'tdesign-vue-next/es/style/index.css'
import './style.css'
import App from './App.vue'

const app = createApp(App)
const pinia = createPinia()
const router = createRouter()

app.use(TDesign)
app.use(pinia)
app.use(router)
app.mount('#app')
