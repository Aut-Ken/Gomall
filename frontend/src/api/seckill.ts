import api from './request';
import { ApiResponse } from './types';

export interface SeckillProduct {
  id: number;
  product_id: number;
  product_name: string;
  product_image: string;
  original_price: number;
  seckill_price: number;
  stock: number;
  start_time: string;
  end_time: string;
}

export interface SeckillResult {
  order_no: string;
  success: boolean;
  message: string;
}

export const seckillApi = {
  seckill: (product_id: number) =>
    api.post<ApiResponse<SeckillResult>>('/seckill', { product_id }),
  initStock: (product_id: number) =>
    api.post<ApiResponse<null>>('/seckill/init', { product_id }),
};
