import { useEffect, useMemo, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';

const navItems = [
  { path: '/chat', label: '聊天' },
  { path: '/contacts', label: '联系人' },
  { path: '/groups', label: '群组' },
];

export default function Layout() {
  const { userInfo, logout } = useAuthStore();
  const { wsState, wsHint } = useChatStore();
  const navigate = useNavigate();
  const location = useLocation();
  const [menuOpen, setMenuOpen] = useState(false);

  const currentPath = useMemo(
    () => navItems.find((item) => location.pathname.startsWith(item.path))?.path || '/chat',
    [location.pathname],
  );

  useEffect(() => {
    const onDocClick = () => setMenuOpen(false);
    document.addEventListener('click', onDocClick);
    return () => document.removeEventListener('click', onDocClick);
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div style={styles.shell}>
      <aside style={styles.sidebar}>
        <div style={styles.brand}>
          <div style={styles.brandMark}>W</div>
          <div>
            <div style={styles.brandTitle}>WsChat</div>
            <div style={styles.brandSub}>即时沟通工作台</div>
          </div>
        </div>

        <nav style={styles.nav}>
          {navItems.map((item) => (
            <button
              key={item.path}
              type="button"
              style={{
                ...styles.navItem,
                ...(currentPath === item.path ? styles.navItemActive : {}),
              }}
              onClick={() => navigate(item.path)}
            >
              {item.label}
            </button>
          ))}
        </nav>

        <div style={styles.statusCard}>
          <div style={styles.statusLabel}>连接状态</div>
          <div style={styles.statusValue}>{wsHint}</div>
          <div style={styles.statusMeta}>{wsState}</div>
        </div>
      </aside>

      <main style={styles.main}>
        <header style={styles.topbar}>
          <div>
            <div style={styles.topTitle}>WsChat</div>
            <div style={styles.topSub}>联系人、会话和群组已接入</div>
          </div>

          <div style={styles.userArea} onClick={(e) => e.stopPropagation()}>
            <div style={styles.userChip}>
              <div style={styles.avatar}>{(userInfo?.nickname || userInfo?.username || 'U').charAt(0)}</div>
              <div style={styles.userMeta}>
                <div style={styles.userName}>{userInfo?.nickname || userInfo?.username || '用户'}</div>
                <div style={styles.userDesc}>在线</div>
              </div>
            </div>
            <button type="button" style={styles.menuBtn} onClick={() => setMenuOpen((v) => !v)}>
              菜单
            </button>
            {menuOpen && (
              <div style={styles.menu}>
                <button type="button" style={styles.menuItem} onClick={handleLogout}>
                  退出登录
                </button>
              </div>
            )}
          </div>
        </header>

        <section style={styles.content}>
          <Outlet />
        </section>
      </main>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  shell: {
    minHeight: '100vh',
    display: 'grid',
    gridTemplateColumns: '280px minmax(0, 1fr)',
    background: 'linear-gradient(180deg, #f5f9ff 0%, #eef4ff 100%)',
  },
  sidebar: {
    padding: 18,
    borderRight: '1px solid #dbe7fb',
    background: '#fff',
    display: 'flex',
    flexDirection: 'column',
    gap: 16,
  },
  brand: { display: 'flex', alignItems: 'center', gap: 12 },
  brandMark: {
    width: 46,
    height: 46,
    borderRadius: 16,
    display: 'grid',
    placeItems: 'center',
    background: '#4a90d9',
    color: '#fff',
    fontWeight: 800,
  },
  brandTitle: { fontSize: 18, fontWeight: 800, color: '#1f2d3d' },
  brandSub: { marginTop: 4, fontSize: 12, color: '#7a869a' },
  nav: { display: 'grid', gap: 10 },
  navItem: {
    height: 46,
    borderRadius: 12,
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
    color: '#456',
    fontWeight: 600,
  },
  navItemActive: {
    background: '#4a90d9',
    borderColor: '#4a90d9',
    color: '#fff',
  },
  statusCard: {
    marginTop: 'auto',
    borderRadius: 16,
    padding: 14,
    background: '#f7fbff',
    border: '1px solid #dbe7fb',
  },
  statusLabel: { fontSize: 12, color: '#7a869a' },
  statusValue: { marginTop: 6, fontSize: 14, fontWeight: 700, color: '#1f2d3d' },
  statusMeta: { marginTop: 4, fontSize: 12, color: '#7a869a', textTransform: 'uppercase' },
  main: { minWidth: 0, display: 'flex', flexDirection: 'column' },
  topbar: {
    height: 64,
    padding: '0 20px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    background: '#4a90d9',
    color: '#fff',
    position: 'relative',
  },
  topTitle: { fontSize: 20, fontWeight: 800 },
  topSub: { marginTop: 4, fontSize: 12, opacity: 0.9 },
  userArea: { display: 'flex', alignItems: 'center', gap: 10, position: 'relative' },
  userChip: {
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    padding: '8px 12px',
    borderRadius: 14,
    background: 'rgba(255,255,255,0.14)',
  },
  avatar: {
    width: 34,
    height: 34,
    borderRadius: 999,
    background: '#fff',
    color: '#4a90d9',
    display: 'grid',
    placeItems: 'center',
    fontWeight: 800,
  },
  userMeta: { minWidth: 0 },
  userName: { fontSize: 13, fontWeight: 700, lineHeight: 1.2 },
  userDesc: { marginTop: 2, fontSize: 11, opacity: 0.9 },
  menuBtn: {
    height: 36,
    padding: '0 12px',
    borderRadius: 10,
    border: '1px solid rgba(255,255,255,0.35)',
    background: 'transparent',
    color: '#fff',
  },
  menu: {
    position: 'absolute',
    right: 0,
    top: 'calc(100% + 10px)',
    minWidth: 140,
    borderRadius: 12,
    padding: 8,
    background: '#fff',
    border: '1px solid #dbe7fb',
    boxShadow: '0 12px 30px rgba(30, 64, 175, 0.14)',
    zIndex: 20,
  },
  menuItem: {
    width: '100%',
    height: 40,
    textAlign: 'left',
    border: 'none',
    borderRadius: 10,
    background: '#fff',
    color: '#1f2d3d',
    padding: '0 12px',
  },
  content: {
    flex: 1,
    minHeight: 0,
  },
};
