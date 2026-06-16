import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '@/api/client';
import AvatarView from '@/components/AvatarView';
import { useMobile } from '@/hooks/useMobile';
import { resolveAvatarUrl } from '@/utils/avatar';
import type { UserInfoResp } from '@/types';

function safeText(value: string | undefined, fallback: string) {
  if (!value) return fallback;
  if (value.includes('�')) return fallback;
  return value;
}

export default function ContactSearch() {
  const navigate = useNavigate();
  const isMobile = useMobile();
  const [keyword, setKeyword] = useState('');
  const [results, setResults] = useState<UserInfoResp[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    const value = keyword.trim();
    if (!value) {
      setResults([]);
      return;
    }
    setLoading(true);
    try {
      const resp = await api.searchUsers(value);
      if (resp.code === 0) {
        setResults(resp.data || []);
      } else {
        alert(resp.message);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleApply = async (toId?: number) => {
    if (!toId) return;
    const resp = await api.applyContact(toId, '你好，想加个好友');
    alert(resp.message);
  };

  return (
    <div style={{ ...styles.wrap, ...(isMobile ? styles.wrapMobile : {}) }}>
      <div style={styles.header}>
        <button type="button" style={styles.backBtn} onClick={() => navigate('/contacts')}>
          返回联系人
        </button>
        <div>
          <div style={styles.title}>查找用户</div>
          <div style={styles.sub}>独立搜索页，不会覆盖联系人列表</div>
        </div>
      </div>

      <div style={styles.searchBar}>
        <input
          style={styles.input}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="搜索用户名或昵称"
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <button type="button" style={styles.searchBtn} onClick={handleSearch} disabled={loading}>
          {loading ? '搜索中' : '搜索'}
        </button>
      </div>

      <div style={styles.list}>
        {results.length === 0 ? (
          <div style={styles.empty}>请输入关键词后搜索</div>
        ) : (
          results.map((user) => {
            const name = safeText(user.nickname || user.username, `用户 ${user.user_id}`);
            return (
              <div key={user.user_id} style={{ ...styles.item, ...(isMobile ? styles.itemMobile : {}) }}>
                <div style={styles.itemBody}>
                  <AvatarView
                    src={resolveAvatarUrl(user.avatar)}
                    alt={name}
                    size={isMobile ? 32 : 36}
                    radius={8}
                    style={styles.avatar}
                  />
                  <div style={styles.itemInfo}>
                    <div style={styles.name}>{name}</div>
                    <div style={styles.subText}>{safeText(user.bio, '暂无签名')}</div>
                  </div>
                </div>
                <button type="button" style={styles.addBtn} onClick={() => handleApply(user.user_id)}>
                  添加好友
                </button>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  wrap: {
    flex: 1,
    minHeight: 0,
    padding: 12,
    display: 'flex',
    flexDirection: 'column',
    gap: 12,
    background: 'linear-gradient(180deg, rgba(245,249,255,1) 0%, rgba(238,244,255,1) 100%)',
  },
  wrapMobile: {
    padding: 8,
    gap: 8,
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    gap: 12,
    padding: 16,
    borderRadius: 18,
    background: 'rgba(255,255,255,0.94)',
    border: '1px solid #dbe7fb',
    boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)',
  },
  backBtn: {
    height: 36,
    padding: '0 12px',
    borderRadius: 10,
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
    color: '#4a90d9',
    fontWeight: 700,
  },
  title: { fontSize: 18, fontWeight: 800, color: '#1f2d3d' },
  sub: { marginTop: 4, fontSize: 12, color: '#7a869a' },
  searchBar: {
    display: 'flex',
    gap: 8,
    padding: 12,
    borderRadius: 18,
    background: 'rgba(255,255,255,0.94)',
    border: '1px solid #dbe7fb',
    boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)',
  },
  input: {
    flex: 1,
    minWidth: 0,
    borderRadius: 12,
    border: '1px solid #dbe7fb',
    padding: '10px 12px',
    fontSize: 14,
    outline: 'none',
  },
  searchBtn: {
    height: 42,
    padding: '0 16px',
    borderRadius: 12,
    border: 'none',
    background: '#4a90d9',
    color: '#fff',
    fontWeight: 700,
  },
  list: {
    flex: 1,
    minHeight: 0,
    overflow: 'auto',
    borderRadius: 18,
    background: 'rgba(255,255,255,0.94)',
    border: '1px solid #dbe7fb',
    boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)',
  },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
  item: {
    display: 'flex',
    alignItems: 'center',
    gap: 12,
    padding: '12px 14px',
    borderBottom: '1px solid #f3f5f7',
  },
  itemMobile: {
    padding: '10px 12px',
  },
  itemBody: {
    flex: 1,
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    minWidth: 0,
  },
  avatar: { background: '#e6eefc' },
  itemInfo: { flex: 1, minWidth: 0 },
  name: { fontSize: 14, fontWeight: 600, color: '#1f2d3d' },
  subText: {
    marginTop: 2,
    fontSize: 12,
    color: '#8a94a6',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  addBtn: {
    height: 34,
    padding: '0 12px',
    borderRadius: 10,
    border: 'none',
    background: '#4a90d9',
    color: '#fff',
    fontWeight: 700,
    flexShrink: 0,
  },
};
