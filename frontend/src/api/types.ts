// API 基础配置
export const API_BASE_URL = '/api';

// 响应类型定义
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 分页参数
export interface PaginationParams {
  page: number;
  page_size: number;
}

// 分页响应
export interface PaginatedResponse<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}
