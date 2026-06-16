import { useEffect, useMemo, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import AvatarView from '@/components/AvatarView';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import { useMobile } from '@/hooks/useMobile';
import { resolveAvatarUrl } from '@/utils/avatar';

const navItems = [
  { path: '/chat', label: '聊天' },
  { path: '/contacts', label: '联系人' },
  { path: '/groups', label: '群组' },
  { path: '/profile', label: '个人信息' },
];

export default function Layout() {
  const { userInfo, logout } = useAuthStore();
  const { wsState, wsHint } = useChatStore();
  const isMobile = useMobile();
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
    <div style={{ ...styles.shell, ...(isMobile ? styles.shellMobile : {}) }}>
      {!isMobile && (
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
            <div style={styles.statusRow}>
              <span
                style={{
                  ...styles.statusDot,
                  ...(wsState === 'connected'
                    ? styles.statusDotOnline
                    : wsState === 'reconnecting' || wsState === 'connecting'
                      ? styles.statusDotBusy
                      : styles.statusDotOffline),
                }}
              />
              <div style={styles.statusValue}>{wsHint}</div>
            </div>
            <div style={styles.statusMeta}>{wsState}</div>
          </div>
        </aside>
      )}

      <main style={{ ...styles.main, ...(isMobile ? styles.mainMobile : {}) }}>
        {!isMobile && (
          <header style={styles.topbar}>
            <div style={styles.topbarInfo}>
              <div style={styles.topTitle}>WsChat</div>
              <div style={styles.topSub}>联系人、会话和群组已接入</div>
            </div>

            <div style={styles.userArea} onClick={(e) => e.stopPropagation()}>
              <div style={styles.userChip}>
                <AvatarView
                  src={resolveAvatarUrl(userInfo?.avatar)}
                  alt={userInfo?.nickname || userInfo?.username || '用户'}
                  size={34}
                  radius={999}
                  style={styles.avatar}
                />
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
                  <button type="button" style={styles.menuItem} onClick={() => navigate('/profile')}>
                    个人信息
                  </button>
                  <button type="button" style={styles.menuItem} onClick={handleLogout}>
                    退出登录
                  </button>
                </div>
              )}
            </div>
          </header>
        )}

        {/* 手机端顶部导航栏 */}
        {isMobile && (
          <nav style={styles.topNav}>
            {navItems.map((item) => (
              <button
                key={item.path}
                type="button"
                style={{
                  ...styles.topNavItem,
                  ...(currentPath === item.path ? styles.topNavItemActive : {}),
                }}
                onClick={() => navigate(item.path)}
              >
                {item.label}
              </button>
            ))}
            <button
              type="button"
              style={{ ...styles.topNavItem, ...styles.logoutBtn }}
              onClick={handleLogout}
            >
              退出
            </button>
          </nav>
        )}

        <section style={{ ...styles.content, ...(isMobile ? styles.contentMobile : {}) }}>
          <Outlet />
        </section>
      </main>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  shell: {
    height: '100dvh',
    display: 'grid',
    gridTemplateColumns: '236px minmax(0, 1fr)',
    background: 'linear-gradient(180deg, #f5f9ff 0%, #eef4ff 100%)',
    overflow: 'hidden',
    width: 'min(1320px, calc(100vw - 72px))',
    margin: '0 auto',
    boxShadow: '0 24px 60px rgba(31, 64, 122, 0.08)',
  },
  shellMobile: {
    width: '100%',
    gridTemplateColumns: 'minmax(0, 1fr)',
    boxShadow: 'none',
    overflow: 'auto',
    WebkitOverflowScrolling: 'touch',
  },
  sidebar: {
    padding: 16,
    borderRight: '1px solid #dbe7fb',
    background: '#fff',
    display: 'flex',
    flexDirection: 'column',
    gap: 14,
  },
  brand: { display: 'flex', alignItems: 'center', gap: 12 },
  brandMark: {
    width: 42,
    height: 42,
    borderRadius: 14,
    display: 'grid',
    placeItems: 'center',
    background: '#4a90d9',
    color: '#fff',
    fontWeight: 800,
  },
  brandTitle: { fontSize: 16, fontWeight: 800, color: '#1f2d3d' },
  brandSub: { marginTop: 4, fontSize: 12, color: '#7a869a' },
  nav: { display: 'grid', gap: 10 },
  navItem: {
    height: 42,
    borderRadius: 10,
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
  statusRow: { marginTop: 6, display: 'flex', alignItems: 'center', gap: 8 },
  statusDot: { width: 10, height: 10, borderRadius: 999, display: 'inline-block' },
  statusDotOnline: { background: '#13c26b', boxShadow: '0 0 0 4px rgba(19,194,107,0.16)' },
  statusDotBusy: { background: '#ffb020', boxShadow: '0 0 0 4px rgba(255,176,32,0.16)' },
  statusDotOffline: { background: '#a0aec0', boxShadow: '0 0 0 4px rgba(160,174,192,0.14)' },
  statusValue: { fontSize: 14, fontWeight: 700, color: '#1f2d3d' },
  statusMeta: { marginTop: 4, fontSize: 12, color: '#7a869a', textTransform: 'uppercase' },
  main: { minWidth: 0, minHeight: 0, display: 'flex', flexDirection: 'column' },
  mainMobile: { width: '100%' },
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
  topbarInfo: { minWidth: 0 },
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
    display: 'flex',
    overflow: 'hidden',
  },
  // 手机端顶部导航（原底部导航改为顶部）
  topNav: {
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    padding: '8px 10px',
    background: '#4a90d9',
    borderBottom: '1px solid #dbe7fb',
    zIndex: 40,
  },
  topNavItem: {
    height: 34,
    padding: '0 14px',
    borderRadius: 10,
    border: '1px solid rgba(255,255,255,0.3)',
    background: 'rgba(255,255,255,0.12)',
    color: '#fff',
    fontSize: 13,
    fontWeight: 700,
    whiteSpace: 'nowrap',
  },
  topNavItemActive: {
    background: '#fff',
    borderColor: '#fff',
    color: '#4a90d9',
  },
  logoutBtn: {
    marginLeft: 'auto',
    background: 'rgba(255,80,80,0.25)',
    borderColor: 'rgba(255,80,80,0.35)',
  },
  // 手机端内容区（顶部留出导航栏空间）
  contentMobile: {
    paddingTop: 0,
  },
};
