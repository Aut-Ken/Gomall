import api from './request';
import { ApiResponse, PaginatedResponse } from './types';

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  stock: number;
  category: string;
  image_url: string;
  status: number;
  created_at: string;
}

export interface ProductListParams {
  page?: number;
  page_size?: number;
  category?: string;
  keyword?: string;
}

export const productApi = {
  getList: (params?: ProductListParams) =>
    api.get<ApiResponse<PaginatedResponse<Product>>>('/product', { params }),
  getDetail: (id: number) =>
    api.get<ApiResponse<Product>>(`/product/${id}`),
  create: (data: Partial<Product>) =>
    api.post<ApiResponse<Product>>('/product', data),
  update: (id: number, data: Partial<Product>) =>
    api.put<ApiResponse<Product>>(`/product/${id}`, data),
  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/product/${id}`),
};
