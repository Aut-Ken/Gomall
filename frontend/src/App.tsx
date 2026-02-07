import { useEffect, useState } from 'react';
import { Routes, Route } from 'react-router-dom';
import { useAuthStore } from './store';
import { userApi } from './api/user';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Products from './pages/Products';
import ProductDetail from './pages/ProductDetail';
import Cart from './pages/Cart';
import Orders from './pages/Orders';
import Seckill from './pages/Seckill';

// 简单的加载骨架屏
// 简单的加载骨架屏
// function Loading() {
//   return (
//     <div style={{
//       display: 'flex',
//       justifyContent: 'center',
//       alignItems: 'center',
//       height: '100vh',
//       background: '#f8fafc'
//     }}>
//       <div style={{
//         width: '40px',
//         height: '40px',
//         border: '3px solid #e2e8f0',
//         borderTopColor: '#2563eb',
//         borderRadius: '50%',
//         animation: 'spin 0.8s linear infinite'
//       }} />
//       <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
//     </div>
//   );
// }

function App() {
  const { setAuth, isAuthenticated } = useAuthStore();

  const [isAuthChecking, setIsAuthChecking] = useState(true);

  console.log('App rendering, isAuthenticated:', isAuthenticated);

  // 初始化时检查登录状态
  useEffect(() => {
    const initAuth = async () => {
      const savedToken = localStorage.getItem('token');
      if (savedToken && !isAuthenticated) {
        try {
          const res = await userApi.getProfile();
          if (res.code === 0 && res.data) {
            setAuth(res.data, savedToken);
          }
        } catch (error) {
          localStorage.removeItem('token');
        }
      }
      setIsAuthChecking(false);
    };

    initAuth();
  }, [isAuthenticated, setAuth]); // Removed token from dependency to prevent infinite loop if token changes

  if (isAuthChecking) {
    return <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>Loading...</div>;
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Header />
      <main style={{ flex: 1 }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/products" element={<Products />} />
          <Route path="/products/:id" element={<ProductDetail />} />
          <Route path="/cart" element={<Cart />} />
          <Route path="/orders" element={<Orders />} />
          <Route path="/seckill" element={<Seckill />} />
          <Route path="*" element={<div style={{ padding: '100px 20px', textAlign: 'center' }}><h1>404 - 页面不存在</h1></div>} />
        </Routes>
      </main>
      <Footer />
    </div>
  );
}

export default App;
