import { useState } from 'react';
import { CartItem } from '../api/cart';
import { useCartStore } from '../store';
import toast from 'react-hot-toast';
import styles from './CartItem.module.css';

interface CartItemProps {
  item: CartItem;
}

export default function CartItemComponent({ item }: CartItemProps) {
  const { updateItem, removeItem } = useCartStore();
  const [loading, setLoading] = useState(false);

  const handleUpdate = async (newQuantity: number) => {
    if (newQuantity < 1) return;
    setLoading(true);
    try {
      await updateItem(item.product_id, newQuantity);
      toast.success('已更新');
    } catch (error: any) {
      toast.error(error.message || '更新失败');
    } finally {
      setLoading(false);
    }
  };

  const handleRemove = async () => {
    if (!window.confirm('确定要从购物车删除此商品吗？')) return;
    setLoading(true);
    try {
      await removeItem(item.product_id);
      toast.success('已删除');
    } catch (error: any) {
      toast.error(error.message || '删除失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.item}>
      <div className={styles.product}>
        <div className={styles.image}>
          {item.product_image ? (
            <img src={item.product_image} alt={item.product_name} />
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
        <div className={styles.info}>
          <h3>{item.product_name}</h3>
          <p className={styles.price}>¥{(item.price || 0).toFixed(2)}</p>
        </div>
      </div>

      <div className={styles.quantity}>
        <button
          onClick={() => handleUpdate(item.quantity - 1)}
          disabled={loading || item.quantity <= 1}
        >
          -
        </button>
        <span>{item.quantity}</span>
        <button
          onClick={() => handleUpdate(item.quantity + 1)}
          disabled={loading}
        >
          +
        </button>
      </div>

      <div className={styles.subtotal}>
        ¥{((item.price || 0) * (item.quantity || 1)).toFixed(2)}
      </div>

      <button className={styles.remove} onClick={handleRemove} disabled={loading}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <polyline points="3,6 5,6 21,6" />
          <path d="M19,6v14a2,2 0 0 1-2,2H7a2,2 0 0 1-2-2V6m3,0V4a2,2 0 0 1 2-2h4a2,2 0 0 1 2,2v2" />
        </svg>
      </button>
    </div>
  );
}
