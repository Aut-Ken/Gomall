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
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // 防止重复点击
    if (isSubmitting) return;
    setIsSubmitting(true);

    if (!formData.username || !formData.password) {
      setIsSubmitting(false);
      toast.error('请填写完整的登录信息');
      return;
    }

    setLoading(true);
    try {
      // 调用登录API，响应拦截器已经返回了 response.data
      const responseData = await userApi.login(formData);

      console.log('登录响应数据:', responseData);

      if (responseData.code === 0) {
        if (responseData.data?.user && responseData.data?.token) {
          setAuth(responseData.data.user, responseData.data.token);
          toast.success('登录成功！');
          navigate('/');
        } else {
          toast.error('登录成功，但数据异常');
        }
      } else {
        toast.error(responseData.message || '登录失败');
      }
    } catch (error: unknown) {
      const err = error as Error;
      console.error('登录异常:', err);
      toast.error(err.message || '登录失败，请稍后重试');
    } finally {
      setIsSubmitting(false);
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

          <button
            type="submit"
            className={`btn btn-primary w-full ${styles.submit}`}
            disabled={loading || isSubmitting}
          >
            {loading || isSubmitting ? '登录中...' : '登录'}
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