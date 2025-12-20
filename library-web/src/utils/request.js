import axios from 'axios';
import { getToken, getRefreshToken, setToken, setRefreshToken, clearAuth } from './auth';
import { $message } from './toast'; // 确保你上一步创建了 toast.js

const service = axios.create({
  baseURL: '', // 配合 vite.config.js 的 proxy 使用
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
});

// 请求拦截器
service.interceptors.request.use(
  config => {
    const token = getToken();
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  error => Promise.reject(error)
);

// 响应拦截器
service.interceptors.response.use(
  response => {
    // 1. 安全检查：如果没有 data，可能是空响应
    if (!response || !response.data) {
      console.warn('后端返回了空数据:', response);
      throw new Error('后端返回数据为空');
    }

    const res = response.data;

    // 2. 安全检查：如果 res 只是一个字符串（例如后端报错返回了 HTML 页面），手动抛错
    if (typeof res === 'string') {
      console.error('后端返回了非 JSON 格式:', res);
      $message.error('系统接口异常，请查看控制台');
      return Promise.reject(new Error('API Response Error'));
    }

    // 3. 正常逻辑：判断业务状态码
    // 注意：这里要防止 res.code 不存在的情况
    if (res.code === 200 || res.code === 201) {
      return res.data; // 返回具体数据
    } else {
      // 业务错误（如密码错误）
      const msg = res.message || '未知业务错误';
      $message.error(msg);
      return Promise.reject(new Error(msg));
    }
  },
  async error => {
    const originalRequest = error.config;
    
    // 处理 401 Token 过期
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const refreshToken = getRefreshToken();
      
      if (refreshToken) {
        try {
          // 注意：这里必须用 axios 原生实例，不能用 service，否则死循环
          const { data } = await axios.post('/api/users/refresh-token', {
            refresh_token: refreshToken
          });

          if (data && data.code === 200) {
            setToken(data.data.access_token);
            setRefreshToken(data.data.refresh_token);
            originalRequest.headers['Authorization'] = `Bearer ${data.data.access_token}`;
            return axios(originalRequest);
          }
        } catch (refreshError) {
          console.error('Token 刷新失败', refreshError);
        }
      }
      clearAuth();
      window.location.href = '/login';
      return Promise.reject(error);
    }
    
    // 处理通用 HTTP 错误
    let msg = error.message;
    if (error.response && error.response.data) {
       // 尝试获取后端返回的详细错误信息
       msg = error.response.data.message || error.response.statusText;
    }
    
    $message.error(msg);
    return Promise.reject(error);
  }
);

export default service;