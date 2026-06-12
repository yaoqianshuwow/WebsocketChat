import { useState } from 'react';
import { useAuthStore } from '@/store/auth';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';

const navItems = [
  { path: '/chat', label: '聊天', icon: '💬' },
  { path: '/contacts', label: '联系人', icon: '👥' },
  { path: '/groups', label: '群组', icon: '🏠' },
];

export default function Layout() {
  const { userInfo, logout } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();
  const [showLogout, setShowLogout] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div style={styles.container}>
      {/* 顶部导航 */}
      <header style={styles.header}>
        <h1 style={styles.title}>WsChat</h1>
        <div style={styles.userArea}>
          <span style={styles.username}>{userInfo?.nickname || userInfo?.username || '用户'}</span>
          <button style={styles.logoutBtn} onClick={() => setShowLogout(!showLogout)}>▼</button>
          {showLogout && (
            <div style={styles.dropdown}>
              <button style={styles.dropdownItem} onClick={handleLogout}>退出登录</button>
            </div>
          )}
        </div>
      </header>

      <div style={styles.body}>
        {/* 侧边导航 */}
        <nav style={styles.sidebar}>
          {navItems.map((item) => (
            <button
              key={item.path}
              style={{
                ...styles.navBtn,
                ...(location.pathname === item.path ? styles.navBtnActive : {}),
              }}
              onClick={() => navigate(item.path)}
            >
              <span style={styles.navIcon}>{item.icon}</span>
              <span>{item.label}</span>
            </button>
          ))}
        </nav>

        {/* 主内容区 */}
        <main style={styles.main}>
          <Outlet />
        </main>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { height: '100vh', display: 'flex', flexDirection: 'column', background: '#f5f5f5' },
  header: {
    display: 'flex', alignItems: 'center', justifyContent: 'space-between',
    padding: '0 20px', height: 56, background: '#4a90d9', color: '#fff',
  },
  title: { margin: 0, fontSize: 20 },
  userArea: { display: 'flex', alignItems: 'center', gap: 8, position: 'relative' },
  username: { fontSize: 14 },
  logoutBtn: { background: 'none', border: 'none', color: '#fff', cursor: 'pointer', fontSize: 12 },
  dropdown: { position: 'absolute', top: 32, right: 0, background: '#fff', borderRadius: 4, boxShadow: '0 2px 8px rgba(0,0,0,0.15)', zIndex: 100 },
  dropdownItem: { display: 'block', width: '100%', padding: '8px 16px', border: 'none', background: 'none', cursor: 'pointer', fontSize: 14, color: '#333', whiteSpace: 'nowrap' },
  body: { display: 'flex', flex: 1, overflow: 'hidden' },
  sidebar: {
    width: 72, background: '#3a7bd5', display: 'flex', flexDirection: 'column',
    padding: '8px 0', gap: 4,
  },
  navBtn: {
    display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 4,
    padding: '12px 0', border: 'none', background: 'none', color: 'rgba(255,255,255,0.7)',
    cursor: 'pointer', fontSize: 12, transition: 'all 0.2s',
  },
  navBtnActive: { background: 'rgba(255,255,255,0.15)', color: '#fff' },
  navIcon: { fontSize: 22 },
  main: { flex: 1, overflow: 'hidden', display: 'flex' },
};
