const TOKEN_KEY = 'library_access_token';
const REFRESH_TOKEN_KEY = 'library_refresh_token';
const USER_KEY = 'library_user_info';

export const setToken = (token) => localStorage.setItem(TOKEN_KEY, token);
export const getToken = () => localStorage.getItem(TOKEN_KEY);
export const removeToken = () => localStorage.removeItem(TOKEN_KEY);

export const setRefreshToken = (token) => localStorage.setItem(REFRESH_TOKEN_KEY, token);
export const getRefreshToken = () => localStorage.getItem(REFRESH_TOKEN_KEY);
export const removeRefreshToken = () => localStorage.removeItem(REFRESH_TOKEN_KEY);

export const setUser = (user) => {
  // 关键修复：如果 user 是 undefined 或 null，绝对不要存进去
  if (user && typeof user === 'object') {
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  }
};

export const getUser = () => {
  const u = localStorage.getItem(USER_KEY);
  // 关键修复：如果读取到的是字符串 "undefined" 或者空，直接返回 null
  if (!u || u === 'undefined' || u === 'null') {
    return null;
  }
  try {
    return JSON.parse(u);
  } catch (e) {
    console.warn('检测到用户数据损坏，已自动重置');
    removeUser(); // 删掉坏数据
    return null;
  }
};

export const removeUser = () => localStorage.removeItem(USER_KEY);

export const clearAuth = () => {
  removeToken();
  removeRefreshToken();
  removeUser();
};