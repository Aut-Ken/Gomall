import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore, useCartStore } from '../store';
import styles from './Header.module.css';

export default function Header() {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, user, logout } = useAuthStore();
  const { items } = useCartStore();
  const [menuOpen, setMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  const cartCount = (items || []).reduce((sum, item) => sum + item.quantity, 0);

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/products?keyword=${encodeURIComponent(searchQuery)}`);
    }
  };

  const isActive = (path: string) => location.pathname === path;

  return (
    <header className={styles.header}>
      <div className={styles.topBar}>
        <div className="container">
          <div className={styles.topContent}>
            <span>欢迎来到 GoMall！</span>
            <div className={styles.topLinks}>
              {isAuthenticated ? (
                <>
                  <span className={styles.userInfo}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                      <circle cx="12" cy="7" r="4" />
                    </svg>
                    {user?.username}
                  </span>
                  {user?.role === 2 && (
                    <Link to="/admin" className={styles.adminLink}>管理后台</Link>
                  )}
                  <button onClick={handleLogout} className={styles.logoutBtn}>退出</button>
                </>
              ) : (
                <>
                  <Link to="/login">登录</Link>
                  <Link to="/register">注册</Link>
                </>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className={styles.mainBar}>
        <div className="container">
          <div className={styles.mainContent}>
            <Link to="/" className={styles.logo}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
                <polyline points="9,22 9,12 15,12 15,22" />
              </svg>
              GoMall
            </Link>

            <nav className={styles.nav}>
              <Link to="/" className={isActive('/') ? styles.active : ''}>首页</Link>
              <Link to="/products" className={isActive('/products') ? styles.active : ''}>商品</Link>
              <Link to="/seckill" className={isActive('/seckill') ? styles.active : ''}>秒杀</Link>
              <Link to="/orders" className={isActive('/orders') ? styles.active : ''}>我的订单</Link>
            </nav>

            <form onSubmit={handleSearch} className={styles.search}>
              <input
                type="text"
                placeholder="搜索商品..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <button type="submit">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="11" cy="11" r="8" />
                  <path d="M21 21l-4.35-4.35" />
                </svg>
              </button>
            </form>

            <Link to="/cart" className={styles.cart}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="9" cy="21" r="1" />
                <circle cx="20" cy="21" r="1" />
                <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6" />
              </svg>
              {cartCount > 0 && <span className={styles.cartBadge}>{cartCount}</span>}
            </Link>

            <button className={styles.menuBtn} onClick={() => setMenuOpen(!menuOpen)}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                {menuOpen ? (
                  <path d="M18 6L6 18M6 6l12 12" />
                ) : (
                  <path d="M3 12h18M3 6h18M3 18h18" />
                )}
              </svg>
            </button>
          </div>
        </div>
      </div>

      {menuOpen && (
        <div className={styles.mobileMenu}>
          <nav>
            <Link to="/" onClick={() => setMenuOpen(false)}>首页</Link>
            <Link to="/products" onClick={() => setMenuOpen(false)}>商品</Link>
            <Link to="/seckill" onClick={() => setMenuOpen(false)}>秒杀</Link>
            <Link to="/orders" onClick={() => setMenuOpen(false)}>我的订单</Link>
            <Link to="/cart" onClick={() => setMenuOpen(false)}>购物车</Link>
            {isAuthenticated ? (
              <button onClick={handleLogout}>退出登录</button>
            ) : (
              <>
                <Link to="/login" onClick={() => setMenuOpen(false)}>登录</Link>
                <Link to="/register" onClick={() => setMenuOpen(false)}>注册</Link>
              </>
            )}
          </nav>
        </div>
      )}
    </header>
  );
}
