import {createApp} from 'vue'
import App from './App.vue'

import vuetify from './plugins/vuetify'
import 'vuetify/styles' // 引入Vuetify的样式文件


const app = createApp(App)
app.use(vuetify)
app.mount('#app')