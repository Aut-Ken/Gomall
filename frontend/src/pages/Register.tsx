import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { userApi } from '../api/user';
import { useAuthStore } from '../store';
import toast from 'react-hot-toast';
import styles from './Auth.module.css';

export default function Register() {
  const navigate = useNavigate();
  const { setAuth } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
    phone: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.username || !formData.password) {
      toast.error('请填写必填信息');
      return;
    }

    if (formData.password !== formData.confirmPassword) {
      toast.error('两次密码输入不一致');
      return;
    }

    if (formData.password.length < 6) {
      toast.error('密码长度至少6位');
      return;
    }

    setLoading(true);
    try {
      const res = await userApi.register({
        username: formData.username,
        password: formData.password,
        email: formData.email || undefined,
        phone: formData.phone || undefined,
      });
      if (res.code === 0) {
        setAuth(res.data.user, res.data.token);
        toast.success('注册成功！');
        navigate('/');
      } else {
        toast.error(res.message || '注册失败');
      }
    } catch (error: any) {
      toast.error(error.message || '注册失败，请稍后重试');
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
          <h1>创建账号</h1>
          <p>加入 GoMall，开启购物之旅</p>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className="input-group">
            <label htmlFor="username">用户名 <span className={styles.required}>*</span></label>
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
            <label htmlFor="password">密码 <span className={styles.required}>*</span></label>
            <input
              type="password"
              id="password"
              name="password"
              placeholder="请输入密码（至少6位）"
              value={formData.password}
              onChange={handleChange}
            />
          </div>

          <div className="input-group">
            <label htmlFor="confirmPassword">确认密码 <span className={styles.required}>*</span></label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              placeholder="请再次输入密码"
              value={formData.confirmPassword}
              onChange={handleChange}
            />
          </div>

          <div className="input-group">
            <label htmlFor="email">邮箱</label>
            <input
              type="email"
              id="email"
              name="email"
              placeholder="请输入邮箱（可选）"
              value={formData.email}
              onChange={handleChange}
            />
          </div>

          <div className="input-group">
            <label htmlFor="phone">手机号</label>
            <input
              type="tel"
              id="phone"
              name="phone"
              placeholder="请输入手机号（可选）"
              value={formData.phone}
              onChange={handleChange}
            />
          </div>

          <button type="submit" className={`btn btn-primary w-full ${styles.submit}`} disabled={loading}>
            {loading ? '注册中...' : '注册'}
          </button>
        </form>

        <div className={styles.footer}>
          <p>
            已有账号？<Link to="/login">立即登录</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
