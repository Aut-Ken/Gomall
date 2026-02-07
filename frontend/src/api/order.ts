import api from './request';
import { ApiResponse } from './types';

export interface Order {
  id: number;
  order_no: string;
  product_id: number;
  product_name: string;
  product_image: string;
  quantity: number;
  total_price: number;
  status: number;
  pay_type: number;
  created_at: string;
}

export interface CreateOrderParams {
  product_id: number;
  quantity: number;
}

export const orderApi = {
  create: (data: CreateOrderParams) =>
    api.post<ApiResponse<Order>>('/order', data),
  checkout: () =>
    api.post<ApiResponse<Order[]>>('/order/checkout'),
  getList: () =>
    api.get<ApiResponse<{ list: Order[]; total: number }>>('/order'),
  getDetail: (order_no: string) =>
    api.get<ApiResponse<Order>>(`/order/${order_no}`),
  pay: (order_no: string, pay_type: number = 1) =>
    api.post<ApiResponse<null>>(`/order/${order_no}/pay`, { pay_type }),
  cancel: (order_no: string) =>
    api.post<ApiResponse<null>>(`/order/${order_no}/cancel`),
};
