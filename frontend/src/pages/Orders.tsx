import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { orderApi, Order } from '../api/order';
import { useAuthStore } from '../store';
import toast from 'react-hot-toast';
import styles from './Orders.module.css';

export default function Orders() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [paying, setPaying] = useState<string | null>(null);

  useEffect(() => {
    if (isAuthenticated) {
      loadOrders();
    } else {
      setLoading(false);
    }
  }, [isAuthenticated]);

  const loadOrders = async () => {
    try {
      const res = await orderApi.getList();
      if (res.code === 0) {
        setOrders(res.data || []);
      }
    } catch (error) {
      console.error('加载订单失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handlePay = async (order_no: string) => {
    setPaying(order_no);
    try {
      const res = await orderApi.pay(order_no);
      if (res.code === 0) {
        toast.success('支付成功！');
        loadOrders();
      } else {
        toast.error(res.message || '支付失败');
      }
    } catch (error: any) {
      toast.error(error.message || '支付失败');
    } finally {
      setPaying(null);
    }
  };

  const handleCancel = async (order_no: string) => {
    if (!window.confirm('确定要取消该订单吗？')) return;
    try {
      const res = await orderApi.cancel(order_no);
      if (res.code === 0) {
        toast.success('订单已取消');
        loadOrders();
      } else {
        toast.error(res.message || '取消失败');
      }
    } catch (error: any) {
      toast.error(error.message || '取消失败');
    }
  };

  const getStatusText = (status: number) => {
    switch (status) {
      case 1:
        return { text: '待支付', class: styles.statusPending };
      case 2:
        return { text: '已支付', class: styles.statusPaid };
      default:
        return { text: '未知状态', class: '' };
    }
  };

  if (!isAuthenticated) {
    return (
      <div className={styles.container}>
        <div className="container">
          <div className={styles.empty}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14,2 14,8 20,8" />
              <line x1="16" y1="13" x2="8" y2="13" />
              <line x1="16" y1="17" x2="8" y2="17" />
              <polyline points="10,9 9,9 8,9" />
            </svg>
            <h3>请先登录</h3>
            <p>登录后查看您的订单</p>
            <button className="btn btn-primary" onClick={() => navigate('/login')}>
              去登录
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className={styles.container}>
        <div className="container">
          <div className="loading"><div className="spinner" /></div>
        </div>
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className={styles.container}>
        <div className="container">
          <div className={styles.empty}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14,2 14,8 20,8" />
              <line x1="16" y1="13" x2="8" y2="13" />
              <line x1="16" y1="17" x2="8" y2="17" />
            </svg>
            <h3>暂无订单</h3>
            <p>快去挑选心仪的商品吧</p>
            <button className="btn btn-primary" onClick={() => navigate('/products')}>
              去购物
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className="container">
        <div className={styles.header}>
          <h1>我的订单</h1>
          <span className={styles.count}>{orders.length} 个订单</span>
        </div>

        <div className={styles.list}>
          {orders.map((order) => {
            const status = getStatusText(order.status);
            return (
              <div key={order.id} className={styles.orderCard}>
                <div className={styles.orderHeader}>
                  <div className={styles.orderInfo}>
                    <span className={styles.orderNo}>订单号: {order.order_no}</span>
                    <span className={styles.orderDate}>{new Date(order.created_at).toLocaleString()}</span>
                  </div>
                  <span className={`${styles.status} ${status.class}`}>
                    {status.text}
                  </span>
                </div>

                <div className={styles.orderBody}>
                  <div className={styles.product}>
                    <div className={styles.productImage}>
                      {order.product_image ? (
                        <img src={order.product_image} alt={order.product_name} />
                      ) : (
                        <div className={styles.placeholder}>
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                            <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                            <circle cx="8.5" cy="8.5" r="1.5" />
                            <polyline points="21,15 16,10 5,21" />
                          </svg>
                        </div>
                      )}
                    </div>
                    <div className={styles.productInfo}>
                      <h3>{order.product_name}</h3>
                      <p>数量: {order.quantity} 件</p>
                    </div>
                  </div>

                  <div className={styles.orderFooter}>
                    <div className={styles.total}>
                      <span>实付金额:</span>
                      <span className={styles.price}>¥{order.total_price.toFixed(2)}</span>
                    </div>
                    <div className={styles.actions}>
                      {order.status === 1 && (
                        <>
                          <button
                            className="btn btn-primary btn-sm"
                            onClick={() => handlePay(order.order_no)}
                            disabled={paying === order.order_no}
                          >
                            {paying === order.order_no ? '支付中...' : '立即支付'}
                          </button>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={() => handleCancel(order.order_no)}
                          >
                            取消订单
                          </button>
                        </>
                      )}
                      {order.status === 2 && (
                        <span className={styles.paidInfo}>订单已完成</span>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
