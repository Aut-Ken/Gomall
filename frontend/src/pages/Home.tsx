import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { productApi } from '../api/product';
import { Product } from '../api/product';
import ProductCard from '../components/ProductCard';
import styles from './Home.module.css';

export default function Home() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    try {
      const res = await productApi.getList({ page: 1, page_size: 8 });
      if (res.code === 0 && res.data.list) {
        setProducts(res.data.list);
      }
    } catch (error) {
      console.error('加载商品失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const categories = [
    { name: '手机数码', icon: '📱', count: 120 },
    { name: '电脑办公', icon: '💻', count: 85 },
    { name: '家用电器', icon: '📺', count: 64 },
    { name: '服饰鞋包', icon: '👗', count: 230 },
    { name: '美妆护肤', icon: '💄', count: 156 },
    { name: '食品生鲜', icon: '🍎', count: 89 },
    { name: '家居家装', icon: '🏠', count: 72 },
    { name: '礼品鲜花', icon: '💐', count: 43 },
  ];

  return (
    <div className={styles.home}>
      {/* Banner */}
      <section className={styles.banner}>
        <div className="container">
          <div className={styles.bannerContent}>
            <div className={styles.bannerText}>
              <span className={styles.tag}>2024 新品上市</span>
              <h1>GoMall 精品电商</h1>
              <p>高并发分布式秒杀系统，极致性能体验</p>
              <Link to="/products" className="btn btn-primary btn-lg">
                立即购物
              </Link>
            </div>
            <div className={styles.bannerImage}>
              <svg viewBox="0 0 200 200" fill="none">
                <circle cx="100" cy="100" r="80" fill="rgba(37, 99, 235, 0.1)" />
                <circle cx="100" cy="100" r="60" fill="rgba(37, 99, 235, 0.15)" />
                <circle cx="100" cy="100" r="40" fill="rgba(37, 99, 235, 0.2)" />
                <path d="M100 60 L100 140 M60 100 L140 100" stroke="var(--primary)" strokeWidth="4" strokeLinecap="round" />
                <circle cx="100" cy="100" r="15" fill="var(--primary)" />
              </svg>
            </div>
          </div>
        </div>
      </section>

      {/* Categories */}
      <section className={styles.categories}>
        <div className="container">
          <div className={styles.categoryGrid}>
            {categories.map((cat) => (
              <Link key={cat.name} to={`/products?category=${encodeURIComponent(cat.name)}`} className={styles.category}>
                <span className={styles.icon}>{cat.icon}</span>
                <span className={styles.name}>{cat.name}</span>
                <span className={styles.count}>{cat.count}件商品</span>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Flash Sale */}
      <section className={styles.flashSale}>
        <div className="container">
          <div className={styles.sectionHeader}>
            <h2>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polygon points="13,2 3,14 12,14 11,22 21,10 12,10 13,2" />
              </svg>
              限时秒杀
            </h2>
            <Link to="/seckill" className={styles.more}>
              查看更多
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </Link>
          </div>
          <div className={styles.saleCard}>
            <div className={styles.saleInfo}>
              <h3>超级秒杀日</h3>
              <p>爆款商品低至1折</p>
              <div className={styles.countdown}>
                <span>距离结束</span>
                <div className={styles.timer}>
                  <span className={styles.time}>02</span>:
                  <span className={styles.time}>35</span>:
                  <span className={styles.time}>48</span>
                </div>
              </div>
              <Link to="/seckill" className="btn btn-primary">
                立即抢购
              </Link>
            </div>
            <div className={styles.saleProducts}>
              <div className={styles.saleProduct}>
                <div className={styles.saleImage}>📱</div>
                <p className={styles.saleName}>iPhone 15 Pro</p>
                <p className={styles.salePrice}>¥5999 <span>¥7999</span></p>
              </div>
              <div className={styles.saleProduct}>
                <div className={styles.saleImage}>🎧</div>
                <p className={styles.saleName}>无线耳机</p>
                <p className={styles.salePrice}>¥299 <span>¥599</span></p>
              </div>
              <div className={styles.saleProduct}>
                <div className={styles.saleImage}>⌚</div>
                <p className={styles.saleName}>智能手表</p>
                <p className={styles.salePrice}>¥499 <span>¥999</span></p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Hot Products */}
      <section className={styles.hotProducts}>
        <div className="container">
          <div className={styles.sectionHeader}>
            <h2>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 2v20M2 12h20" />
              </svg>
              热门商品
            </h2>
            <Link to="/products" className={styles.more}>
              查看更多
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </Link>
          </div>

          {loading ? (
            <div className="loading"><div className="spinner" /></div>
          ) : (
            <div className="grid grid-4">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>
          )}
        </div>
      </section>

      {/* Features */}
      <section className={styles.features}>
        <div className="container">
          <div className={styles.featureGrid}>
            <div className={styles.feature}>
              <div className={styles.featureIcon}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="1" y="3" width="15" height="13" />
                  <polygon points="16,8 20,8 23,11 23,16 16,16" />
                  <circle cx="5.5" cy="18.5" r="2.5" />
                  <circle cx="18.5" cy="18.5" r="2.5" />
                </svg>
              </div>
              <h3>免费配送</h3>
              <p>订单满99元免运费</p>
            </div>
            <div className={styles.feature}>
              <div className={styles.featureIcon}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
                </svg>
              </div>
              <h3>品质保障</h3>
              <p>7天无理由退换货</p>
            </div>
            <div className={styles.feature}>
              <div className={styles.featureIcon}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="12" cy="12" r="10" />
                  <polyline points="12,6 12,12 16,14" />
                </svg>
              </div>
              <h3>快速响应</h3>
              <p>24小时客户服务</p>
            </div>
            <div className={styles.feature}>
              <div className={styles.featureIcon}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                  <polyline points="3.27,6.96 12,12.01 20.73,6.96" />
                  <line x1="12" y1="22.08" x2="12" y2="12" />
                </svg>
              </div>
              <h3>正品保证</h3>
              <p>官方授权 正品直营</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
