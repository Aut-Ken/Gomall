import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { cartApi } from '../api/cart';
import { CartItem } from '../api/cart';
import { useCartStore, useAuthStore } from '../store';
import CartItemComponent from '../components/CartItem';
import toast from 'react-hot-toast';
import styles from './Cart.module.css';

export default function Cart() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const { items, getTotal, clearCart, setCart } = useCartStore();
  const [loading, setLoading] = useState(true);
  const [clearing, setClearing] = useState(false);

  useEffect(() => {
    loadCart();
  }, []);

  const loadCart = async () => {
    if (!isAuthenticated) {
      setLoading(false);
      return;
    }
    try {
      const res = await cartApi.getList();
      if (res.code === 0) {
        setCart(res.data || []);
      }
    } catch (error) {
      console.error('加载购物车失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleClearCart = async () => {
    if (!window.confirm('确定要清空购物车吗？')) return;
    setClearing(true);
    try {
      const res = await cartApi.clear();
      if (res.code === 0) {
        clearCart();
        toast.success('购物车已清空');
      }
    } catch (error: any) {
      toast.error(error.message || '清空失败');
    } finally {
      setClearing(false);
    }
  };

  const handleCheckout = () => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      navigate('/login');
      return;
    }
    if (items.length === 0) {
      toast.error('购物车是空的');
      return;
    }
    // 跳转到结算页面或创建订单
    navigate('/orders');
  };

  if (loading) {
    return (
      <div className={styles.container}>
        <div className="container">
          <div className="loading"><div className="spinner" /></div>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className={styles.container}>
        <div className="container">
          <div className={styles.empty}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <circle cx="9" cy="21" r="1" />
              <circle cx="20" cy="21" r="1" />
              <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6" />
            </svg>
            <h3>请先登录</h3>
            <p>登录后查看您的购物车</p>
            <Link to="/login" className="btn btn-primary">去登录</Link>
          </div>
        </div>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className={styles.container}>
        <div className="container">
          <div className={styles.empty}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <circle cx="9" cy="21" r="1" />
              <circle cx="20" cy="21" r="1" />
              <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6" />
            </svg>
            <h3>购物车是空的</h3>
            <p>快去挑选心仪的商品吧</p>
            <Link to="/products" className="btn btn-primary">去购物</Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className="container">
        <div className={styles.header}>
          <h1>购物车</h1>
          <span className={styles.count}>{items.length} 件商品</span>
        </div>

        <div className={styles.content}>
          <div className={styles.items}>
            {items.map((item) => (
              <CartItemComponent key={item.product_id} item={item} />
            ))}
          </div>

          <div className={styles.summary}>
            <div className={styles.summaryCard}>
              <h3>订单摘要</h3>

              <div className={styles.summaryRow}>
                <span>商品数量</span>
                <span>{items.reduce((sum, item) => sum + item.quantity, 0)} 件</span>
              </div>

              <div className={styles.summaryRow}>
                <span>商品总价</span>
                <span>¥{getTotal().toFixed(2)}</span>
              </div>

              <div className={styles.summaryRow}>
                <span>运费</span>
                <span className={styles.free}>免运费</span>
              </div>

              <div className={`${styles.summaryRow} ${styles.total}`}>
                <span>应付金额</span>
                <span className={styles.totalPrice}>¥{getTotal().toFixed(2)}</span>
              </div>

              <button className="btn btn-primary w-full" onClick={handleCheckout}>
                去结算
              </button>

              <button
                className={`btn btn-outline w-full ${styles.clearBtn}`}
                onClick={handleClearCart}
                disabled={clearing}
              >
                清空购物车
              </button>
            </div>

            <Link to="/products" className={styles.continue}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M19 12H5M12 19l-7-7 7-7" />
              </svg>
              继续购物
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
