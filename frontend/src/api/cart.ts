import api from './request';
import { ApiResponse } from './types';

export interface CartItem {
  id: number;
  product_id: number;
  product_name: string;
  price: number;
  product_image: string;
  quantity: number;
  sub_total?: number;
}

export interface AddCartParams {
  product_id: number;
  quantity: number;
}

export interface UpdateCartParams {
  product_id: number;
  quantity: number;
}

export const cartApi = {
  getList: () =>
    api.get<ApiResponse<{ items: CartItem[]; total_count: number; total_price: number }>>('/cart'),
  add: (data: AddCartParams) =>
    api.post<ApiResponse<CartItem>>('/cart', data),
  update: (data: UpdateCartParams) =>
    api.put<ApiResponse<CartItem>>('/cart', data),
  remove: (product_id: number) =>
    api.delete<ApiResponse<null>>('/cart', { data: { product_id } }),
  clear: () =>
    api.delete<ApiResponse<null>>('/cart/clear'),
};
