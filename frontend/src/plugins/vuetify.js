// src/plugins/vuetify.js
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css' // 确保安装了图标库

export default createVuetify({
    components,
    directives,
    theme: {
        defaultTheme: 'light', // 或者 'dark'，根据你的需求设置默认主题
        themes: {
            light: {
                primary: '#1867C0', // 自定义主题颜色等
                secondary: '#5CBBF6', // 更多颜色自定义...
            },
            dark: { /* 暗色主题设置 */ },
        },
    },
})