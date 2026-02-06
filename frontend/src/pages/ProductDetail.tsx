import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productApi, Product } from '../api/product';
import { cartApi } from '../api/cart';
import { useAuthStore, useCartStore } from '../store';
import toast from 'react-hot-toast';
import styles from './ProductDetail.module.css';

export default function ProductDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const { addItem, setCart } = useCartStore();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const [addingToCart, setAddingToCart] = useState(false);

  useEffect(() => {
    loadProduct();
  }, [id]);

  const loadProduct = async () => {
    try {
      const res = await productApi.getDetail(Number(id));
      if (res.code === 0 && res.data) {
        setProduct(res.data);
      } else {
        toast.error('商品不存在');
        navigate('/products');
      }
    } catch (error) {
      console.error('加载商品失败:', error);
      toast.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      navigate('/login');
      return;
    }

    if (!product) return;

    setAddingToCart(true);
    try {
      const res = await cartApi.add({ product_id: product.id, quantity });
      if (res.code === 0) {
        // 同步本地购物车状态
        const cartRes = await cartApi.getList();
        if (cartRes.code === 0) {
          setCart(cartRes.data || []);
        }
        toast.success('已添加到购物车');
      } else {
        toast.error(res.message || '添加失败');
      }
    } catch (error: any) {
      toast.error(error.message || '添加失败');
    } finally {
      setAddingToCart(false);
    }
  };

  const handleBuyNow = () => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      navigate('/login');
      return;
    }
    // 直接购买逻辑
    navigate(`/orders?product_id=${id}&quantity=${quantity}`);
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

  if (!product) {
    return (
      <div className={styles.container}>
        <div className="container">
          <div className="empty-state">
            <h3>商品不存在</h3>
            <button className="btn btn-primary" onClick={() => navigate('/products')}>
              返回商品列表
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className="container">
        <div className={styles.breadcrumb}>
          <span onClick={() => navigate('/')}>首页</span>
          <span>/</span>
          <span onClick={() => navigate('/products')}>商品</span>
          <span>/</span>
          <span>{product.name}</span>
        </div>

        <div className={styles.content}>
          <div className={styles.gallery}>
            <div className={styles.mainImage}>
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
            </div>
          </div>

          <div className={styles.info}>
            <h1 className={styles.name}>{product.name}</h1>
            <p className={styles.description}>{product.description}</p>

            <div className={styles.priceSection}>
              <span className={styles.priceLabel}>价格</span>
              <span className={styles.price}>
                <span className={styles.currency}>¥</span>
                <span className={styles.value}>{product.price.toFixed(2)}</span>
              </span>
            </div>

            <div className={styles.meta}>
              <div className={styles.metaItem}>
                <span className={styles.metaLabel}>库存</span>
                <span className={styles.metaValue}>{product.stock} 件</span>
              </div>
              <div className={styles.metaItem}>
                <span className={styles.metaLabel}>分类</span>
                <span className={styles.metaValue}>{product.category || '未分类'}</span>
              </div>
            </div>

            <div className={styles.quantity}>
              <span className={styles.quantityLabel}>数量</span>
              <div className={styles.quantityControl}>
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  disabled={quantity <= 1}
                >
                  -
                </button>
                <span>{quantity}</span>
                <button
                  onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                  disabled={quantity >= product.stock}
                >
                  +
                </button>
              </div>
            </div>

            <div className={styles.actions}>
              <button
                className="btn btn-primary btn-lg"
                onClick={handleAddToCart}
                disabled={addingToCart || product.stock <= 0}
              >
                {addingToCart ? '添加中...' : product.stock <= 0 ? '已售罄' : '加入购物车'}
              </button>
              <button
                className="btn btn-success btn-lg"
                onClick={handleBuyNow}
                disabled={product.stock <= 0}
              >
                立即购买
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
