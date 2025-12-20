import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    // 开启代理，解决跨域问题
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 【重要】这里改成你后端实际运行的地址和端口
        changeOrigin: true,
        // 如果你的后端路由本身就有 /api 前缀（看日志是有的），就不需要 pathRewrite
      }
    }
  }
})