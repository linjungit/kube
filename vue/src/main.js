import { createApp } from 'vue'
import App from './App.vue'
//引入 router
import router from './router'
//引入 ant
import Antd from 'ant-design-vue';
//暗黑风格主题
import 'ant-design-vue/dist/antd.dark.css'
import * as Icons from '@ant-design/icons-vue'

const app = createApp(App)
//图标注册全局组建
for (const i in Icons){
    app.component(i,Icons[i])
}
app.use(Antd)
app.use(router)
app.mount('#app')
