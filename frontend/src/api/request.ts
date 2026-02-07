import axios, { AxiosRequestConfig } from 'axios';
import { API_BASE_URL } from './types';

const instance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器 - 添加token
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理错误
instance.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error.response?.data || error);
  }
);

// 封装请求方法，确保返回类型正确
const api = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) => {
    return instance.get<any, T>(url, config);
  },
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => {
    return instance.post<any, T>(url, data, config);
  },
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => {
    return instance.put<any, T>(url, data, config);
  },
  delete: <T = any>(url: string, config?: AxiosRequestConfig) => {
    return instance.delete<any, T>(url, config);
  },
};

export default api;
