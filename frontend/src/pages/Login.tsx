import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { userApi } from '../api/user';
import { useAuthStore } from '../store';
import toast from 'react-hot-toast';
import styles from './Auth.module.css';

export default function Login() {
  const navigate = useNavigate();
  const { setAuth } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.username || !formData.password) {
      toast.error('请填写完整的登录信息');
      return;
    }

    setLoading(true);
    try {
      const res = await userApi.login(formData);
      if (res.code === 0) {
        setAuth(res.data.user, res.data.token);
        toast.success('登录成功！');
        navigate('/');
      } else {
        toast.error(res.message || '登录失败');
      }
    } catch (error: any) {
      toast.error(error.message || '登录失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <div className={styles.header}>
          <h1>欢迎回来</h1>
          <p>登录 GoMall 账号，享受更多优惠</p>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className="input-group">
            <label htmlFor="username">用户名</label>
            <input
              type="text"
              id="username"
              name="username"
              placeholder="请输入用户名"
              value={formData.username}
              onChange={handleChange}
            />
          </div>

          <div className="input-group">
            <label htmlFor="password">密码</label>
            <input
              type="password"
              id="password"
              name="password"
              placeholder="请输入密码"
              value={formData.password}
              onChange={handleChange}
            />
          </div>

          <button type="submit" className={`btn btn-primary w-full ${styles.submit}`} disabled={loading}>
            {loading ? '登录中...' : '登录'}
          </button>
        </form>

        <div className={styles.footer}>
          <p>
            还没有账号？<Link to="/register">立即注册</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
