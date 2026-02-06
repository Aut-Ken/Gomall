import { Link } from 'react-router-dom';
import { Product } from '../api/product';
import styles from './ProductCard.module.css';

interface ProductCardProps {
  product: Product;
  showAdmin?: boolean;
  onDelete?: (id: number) => void;
}

export default function ProductCard({ product, showAdmin, onDelete }: ProductCardProps) {
  return (
    <div className={styles.card}>
      <Link to={`/products/${product.id}`} className={styles.image}>
        {product.image_url ? (
          <img src={product.image_url} alt={product.name} />
        ) : (
          <div className={styles.placeholder}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
              <circle cx="8.5" cy="8.5" r="1.5" />
              <polyline points="21,15 16,10 5,21" />
            </svg>
          </div>
        )}
        {product.stock <= 0 && (
          <span className={styles.soldOut}>售罄</span>
        )}
      </Link>

      <div className={styles.content}>
        <Link to={`/products/${product.id}`} className={styles.name}>
          {product.name}
        </Link>
        <p className={styles.desc}>{product.description}</p>
        <div className={styles.footer}>
          <div className={styles.price}>
            <span className={styles.currency}>¥</span>
            <span className={styles.value}>{product.price.toFixed(2)}</span>
          </div>
          <span className={styles.stock}>库存: {product.stock}</span>
        </div>
        <div className={styles.actions}>
          <Link to={`/products/${product.id}`} className="btn btn-primary btn-sm">
            查看详情
          </Link>
          {showAdmin && (
            <button className="btn btn-danger btn-sm" onClick={() => onDelete?.(product.id)}>
              删除
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
