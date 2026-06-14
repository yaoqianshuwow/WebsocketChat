import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login, loading } = useAuthStore();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!username || !password) {
      setError('请填写用户名和密码');
      return;
    }

    const ok = await login(username, password);
    if (ok) {
      navigate('/chat');
    } else {
      setError('登录失败，请检查用户名和密码');
    }
  };

  return (
    <div style={styles.container}>
      <form style={styles.form} onSubmit={handleSubmit}>
        <h2 style={styles.title}>WsChat 登录</h2>

        {error && <div style={styles.error}>{error}</div>}

        <div style={styles.field}>
          <label style={styles.label}>用户名</label>
          <input
            style={styles.input}
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="请输入用户名"
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>密码</label>
          <input
            style={styles.input}
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="请输入密码"
          />
        </div>

        <button style={styles.btn} type="submit" disabled={loading}>
          {loading ? '登录中...' : '登录'}
        </button>

        <div style={styles.footer}>
          没有账号？<Link to="/register" style={styles.link}>立即注册</Link>
        </div>
      </form>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center',
    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  },
  form: {
    background: '#fff', padding: 40, borderRadius: 12, boxShadow: '0 10px 40px rgba(0,0,0,0.2)',
    width: 360,
  },
  title: { margin: '0 0 24px', textAlign: 'center', color: '#333', fontSize: 24 },
  error: { background: '#fff2f0', border: '1px solid #ffccc7', color: '#ff4d4f', padding: '8px 12px', borderRadius: 4, marginBottom: 16, fontSize: 13 },
  field: { marginBottom: 16 },
  label: { display: 'block', marginBottom: 6, fontSize: 13, color: '#666' },
  input: {
    width: '100%', padding: '10px 12px', borderRadius: 6, border: '1px solid #e0e0e0',
    fontSize: 14, outline: 'none', boxSizing: 'border-box',
  },
  btn: {
    width: '100%', padding: '12px', background: '#4a90d9', color: '#fff', border: 'none',
    borderRadius: 6, fontSize: 16, cursor: 'pointer', marginTop: 8,
  },
  footer: { textAlign: 'center', marginTop: 16, fontSize: 13, color: '#999' },
  link: { color: '#4a90d9', textDecoration: 'none' },
};
