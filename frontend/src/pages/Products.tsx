import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { productApi, Product } from '../api/product';
import ProductCard from '../components/ProductCard';
import styles from './Products.module.css';

export default function Products() {
  const [searchParams] = useSearchParams();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({
    page: 1,
    page_size: 12,
    total: 0,
  });
  const [filters, setFilters] = useState({
    category: searchParams.get('category') || '',
    keyword: searchParams.get('keyword') || '',
  });

  const categories = ['全部', '手机数码', '电脑办公', '家用电器', '服饰鞋包', '美妆护肤', '食品生鲜', '家居家装', '礼品鲜花'];

  useEffect(() => {
    loadProducts();
  }, [pagination.page, filters.category, filters.keyword]);

  const loadProducts = async () => {
    setLoading(true);
    try {
      const res = await productApi.getList({
        page: pagination.page,
        page_size: pagination.page_size,
        category: filters.category === '全部' ? '' : filters.category,
        keyword: filters.keyword,
      });
      if (res.code === 0 && res.data.list) {
        setProducts(res.data.list);
        setPagination((prev) => ({
          ...prev,
          total: res.data.total,
        }));
      }
    } catch (error) {
      console.error('加载商品失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCategoryChange = (category: string) => {
    setFilters({ ...filters, category });
    setPagination((prev) => ({ ...prev, page: 1 }));
  };

  const handleSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const keyword = formData.get('keyword') as string;
    setFilters({ ...filters, keyword });
    setPagination((prev) => ({ ...prev, page: 1 }));
  };

  const handlePageChange = (page: number) => {
    setPagination((prev) => ({ ...prev, page }));
  };

  const totalPages = Math.ceil(pagination.total / pagination.page_size);

  return (
    <div className={styles.container}>
      <div className="container">
        <div className={styles.header}>
          <div>
            <h1>商品列表</h1>
            <p>共 {pagination.total} 件商品</p>
          </div>
          <form onSubmit={handleSearch} className={styles.search}>
            <input
              type="text"
              name="keyword"
              placeholder="搜索商品..."
              defaultValue={filters.keyword}
            />
            <button type="submit">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="11" cy="11" r="8" />
                <path d="M21 21l-4.35-4.35" />
              </svg>
            </button>
          </form>
        </div>

        <div className={styles.content}>
          <aside className={styles.sidebar}>
            <h3>分类</h3>
            <ul>
              {categories.map((cat) => (
                <li
                  key={cat}
                  className={filters.category === (cat === '全部' ? '' : cat) || (cat === '全部' && !filters.category) ? styles.active : ''}
                  onClick={() => handleCategoryChange(cat === '全部' ? '' : cat)}
                >
                  {cat}
                </li>
              ))}
            </ul>
          </aside>

          <main className={styles.main}>
            {loading ? (
              <div className="loading"><div className="spinner" /></div>
            ) : products.length === 0 ? (
              <div className="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <circle cx="11" cy="11" r="8" />
                  <path d="M21 21l-4.35-4.35" />
                </svg>
                <h3>没有找到商品</h3>
                <p>试试其他关键词或分类</p>
              </div>
            ) : (
              <>
                <div className="grid grid-3">
                  {products.map((product) => (
                    <ProductCard key={product.id} product={product} />
                  ))}
                </div>

                {totalPages > 1 && (
                  <div className={styles.pagination}>
                    <button
                      className={styles.pageBtn}
                      disabled={pagination.page <= 1}
                      onClick={() => handlePageChange(pagination.page - 1)}
                    >
                      上一页
                    </button>
                    <span className={styles.pageInfo}>
                      第 {pagination.page} / {totalPages} 页
                    </span>
                    <button
                      className={styles.pageBtn}
                      disabled={pagination.page >= totalPages}
                      onClick={() => handlePageChange(pagination.page + 1)}
                    >
                      下一页
                    </button>
                  </div>
                )}
              </>
            )}
          </main>
        </div>
      </div>
    </div>
  );
}
