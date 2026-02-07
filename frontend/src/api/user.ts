import api from './request';
import { ApiResponse } from './types';

// 用户相关
export interface User {
  id: number;
  username: string;
  email: string;
  phone: string;
  role: number;
  created_at: string;
}

export interface LoginParams {
  username: string;
  password: string;
}

export interface RegisterParams {
  username: string;
  password: string;
  email?: string;
  phone?: string;
}

export interface ChangePasswordParams {
  old_password: string;
  new_password: string;
}

export const userApi = {
  login: (data: LoginParams) => api.post<ApiResponse<{ token: string; user: User }>>('/user/login', data),
  register: (data: RegisterParams) => api.post<ApiResponse<{ token: string; user: User }>>('/user/register', data),
  getProfile: () => api.get<ApiResponse<User>>('/user/profile'),
  changePassword: (data: ChangePasswordParams) => api.post<ApiResponse<null>>('/auth/change-password', data),
  logout: () => api.post<ApiResponse<null>>('/auth/logout'),
  refreshToken: () => api.post<ApiResponse<{ token: string }>>('/auth/refresh-token'),
};
