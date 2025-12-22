import request from '../utils/request';

// --- 用户模块 ---
export const login = (data) => request.post('/api/users/login', data);
export const register = (data) => request.post('/api/users/register', data);
export const logout = () => request.post('/api/users/logout');
export const getUserInfo = () => request.get('/api/users/me');
export const updateUserInfo = (data) => request.put('/api/users/me', data);
export const changePassword = (data) => request.post('/api/users/change-password', data);

// 管理员用户接口
export const getUserList = (params) => request.get('/api/users', { params });
export const createUser = (data) => request.post('/api/users', data);
export const updateUser = (id, data) => request.put(`/api/users/${id}`, data);
export const deleteUser = (id) => request.delete(`/api/users/${id}`);

// --- 图书模块 ---
export const getBooks = (params) => request.get('/api/books', { params });
export const getBookDetail = (id) => request.get(`/api/books/${id}`);
export const addBook = (data) => request.post('/api/books', data);
export const updateBook = (id, data) => request.put(`/api/books/${id}`, data);
export const deleteBook = (id) => request.delete(`/api/books/${id}`);
export const getPopularBooks = (params) => request.get('/api/books/popular', { params });

// --- 借阅模块 ---
export const borrowBook = (data) => request.post('/api/borrow', data);
export const returnBook = (id, data) => request.post(`/api/borrow/${id}/return`, data);
export const renewBook = (id, data) => request.post(`/api/borrow/${id}/renew`, data);
export const getBorrowRecords = (params) => request.get('/api/borrow', { params });
export const getCurrentBorrows = () => request.get('/api/borrow/current');

// --- 预约模块 ---
export const reserveBook = (data) => request.post('/api/reservations', data);
export const cancelReservation = (id) => request.delete(`/api/reservations/${id}`);
export const getMyReservations = () => request.get('/api/reservations/my');

// --- 分类模块 ---
export const getCategories = (params) => request.get('/api/categories', { params });
export const getCategoryDetail = (id) => request.get(`/api/categories/${id}`);
export const addCategory = (data) => request.post('/api/categories', data);
export const updateCategory = (id, data) => request.put(`/api/categories/${id}`, data);
export const deleteCategory = (id) => request.delete(`/api/categories/${id}`);

// --- 统计模块 ---
export const getStatsOverview = () => request.get('/api/stats/overview');
export const getBorrowStats = (params) => request.get('/api/stats/borrow', { params });
export const getUserStats = (userId) => request.get(`/api/stats/user/${userId}`);
export const getPopularBooksStats = (params) => request.get('/api/stats/popular-books', { params });
export const getCategoryStats = () => request.get('/api/stats/categories');