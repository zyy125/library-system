# 图书管理系统前端

基于 Vue 3 + Vite 开发的图书管理系统前端应用。

## 技术栈

- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **样式**: 原生 CSS (现代化设计)

## 功能特性

### 普通用户功能
- ✅ 用户注册/登录
- ✅ 图书浏览（搜索、筛选、详情）
- ✅ 图书借阅/预约
- ✅ 当前借阅管理（还书、续借）
- ✅ 借阅历史查看
- ✅ 我的预约管理
- ✅ 个人信息编辑
- ✅ 修改密码

### 管理员功能
- ✅ 数据统计看板
- ✅ 用户管理（CRUD）
- ✅ 图书管理（CRUD + 批量导入）
- ✅ 借阅管理（查看所有记录、处理归还）
- ✅ 分类管理（CRUD）

## 快速开始

### 安装依赖

```bash
cd library-web
npm install
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:5173

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 后端配置

后端 API 地址在 `vite.config.js` 中配置：

```javascript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 后端地址
        changeOrigin: true
      }
    }
  }
})
```

## API 文档

详见项目根目录的 `图书管理系统 API 文档.md`

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88
