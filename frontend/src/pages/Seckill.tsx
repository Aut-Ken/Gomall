import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store';
import toast from 'react-hot-toast';
import styles from './Seckill.module.css';

interface SeckillProduct {
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

export default function Seckill() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const [products, setProducts] = useState<SeckillProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [seckilling, setSeckilling] = useState<number | null>(null);
  const [timeLeft, setTimeLeft] = useState({ hours: 0, minutes: 0, seconds: 0 });

  // 模拟秒杀商品数据（实际应从API获取）
  useEffect(() => {
    // 实际项目中这里应该调用秒杀商品列表API
    setProducts([
      {
        id: 1,
        product_id: 1,
        product_name: 'iPhone 15 Pro Max',
        product_image: '',
        original_price: 9999,
        seckill_price: 5999,
        stock: 10,
        start_time: new Date().toISOString(),
        end_time: new Date(Date.now() + 3 * 60 * 60 * 1000).toISOString(),
      },
      {
        id: 2,
        product_id: 2,
        product_name: 'MacBook Pro 14"',
        product_image: '',
        original_price: 14999,
        seckill_price: 9999,
        stock: 5,
        start_time: new Date().toISOString(),
        end_time: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(),
      },
      {
        id: 3,
        product_id: 3,
        product_name: 'AirPods Pro 2',
        product_image: '',
        original_price: 1899,
        seckill_price: 999,
        stock: 50,
        start_time: new Date().toISOString(),
        end_time: new Date(Date.now() + 4 * 60 * 60 * 1000).toISOString(),
      },
      {
        id: 4,
        product_id: 4,
        product_name: 'Apple Watch Ultra 2',
        product_image: '',
        original_price: 6499,
        seckill_price: 4999,
        stock: 20,
        start_time: new Date().toISOString(),
        end_time: new Date(Date.now() + 5 * 60 * 60 * 1000).toISOString(),
      },
    ]);
    setLoading(false);
  }, []);

  // 倒计时
  useEffect(() => {
    const endTime = new Date(Date.now() + 3 * 60 * 60 * 1000).getTime();

    const timer = setInterval(() => {
      const now = new Date().getTime();
      const distance = endTime - now;

      if (distance > 0) {
        setTimeLeft({
          hours: Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
          minutes: Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60)),
          seconds: Math.floor((distance % (1000 * 60)) / 1000),
        });
      }
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const handleSeckill = async (product: SeckillProduct) => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      navigate('/login');
      return;
    }

    if (product.stock <= 0) {
      toast.error('商品已售罄');
      return;
    }

    setSeckilling(product.id);
    try {
      // 实际项目中这里调用秒杀API
      await new Promise((resolve) => setTimeout(resolve, 1000));

      toast.success('秒杀成功！订单已创建');
      navigate('/orders');
    } catch (error: any) {
      toast.error(error.message || '秒杀失败');
    } finally {
      setSeckilling(null);
    }
  };

  return (
    <div className={styles.container}>
      <div className="container">
        {/* 秒杀Banner */}
        <div className={styles.banner}>
          <div className={styles.bannerContent}>
            <h1>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polygon points="13,2 3,14 12,14 11,22 21,10 12,10 13,2" />
              </svg>
              限时秒杀
            </h1>
            <p>爆款商品限时抢购，超值优惠不容错过</p>
            <div className={styles.countdown}>
              <span>本场结束</span>
              <div className={styles.timer}>
                <span className={styles.timeBlock}>
                  {String(timeLeft.hours).padStart(2, '0')}
                </span>
                <span className={styles.separator}>:</span>
                <span className={styles.timeBlock}>
                  {String(timeLeft.minutes).padStart(2, '0')}
                </span>
                <span className={styles.separator}>:</span>
                <span className={styles.timeBlock}>
                  {String(timeLeft.seconds).padStart(2, '0')}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* 秒杀商品列表 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>正在秒杀</h2>

          {loading ? (
            <div className="loading"><div className="spinner" /></div>
          ) : (
            <div className={styles.grid}>
              {products.map((product) => (
                <div key={product.id} className={styles.card}>
                  <div className={styles.cardImage}>
                    {product.product_image ? (
                      <img src={product.product_image} alt={product.product_name} />
                    ) : (
                      <div className={styles.placeholder}>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                          <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                          <circle cx="8.5" cy="8.5" r="1.5" />
                          <polyline points="21,15 16,10 5,21" />
                        </svg>
                      </div>
                    )}
                    <div className={styles.discount}>
                      {Math.round((1 - product.seckill_price / product.original_price) * 100)}%OFF
                    </div>
                  </div>

                  <div className={styles.cardContent}>
                    <h3>{product.product_name}</h3>

                    <div className={styles.progress}>
                      <div
                        className={styles.progressBar}
                        style={{
                          width: `${Math.min(100, ((100 - product.stock) / 100) * 100)}%`,
                        }}
                      />
                      <span className={styles.progressText}>
                        已抢 {100 - product.stock}%
                      </span>
                    </div>

                    <div className={styles.prices}>
                      <div className={styles.seckillPrice}>
                        <span className={styles.label}>秒杀价</span>
                        <span className={styles.currency}>¥</span>
                        <span className={styles.value}>{product.seckill_price}</span>
                      </div>
                      <div className={styles.originalPrice}>
                        ¥{product.original_price}
                      </div>
                    </div>

                    <div className={styles.stock}>
                      仅剩 {product.stock} 件
                    </div>

                    <button
                      className={`btn ${product.stock > 0 ? 'btn-primary' : 'btn-secondary'} w-full`}
                      onClick={() => handleSeckill(product)}
                      disabled={seckilling === product.id || product.stock <= 0}
                    >
                      {seckilling === product.id
                        ? '正在秒杀...'
                        : product.stock <= 0
                        ? '已售罄'
                        : '立即秒杀'}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* 秒杀规则 */}
        <div className={styles.rules}>
          <h3>秒杀规则</h3>
          <ul>
            <li>秒杀商品数量有限，售完即止</li>
            <li>每个用户限秒杀1件商品</li>
            <li>秒杀成功后需在30分钟内完成支付</li>
            <li>秒杀商品不支持退换货</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
