import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '@/api/client';
import AvatarView from '@/components/AvatarView';
import { useChatStore } from '@/store/chat';
import { resolveAvatarUrl } from '@/utils/avatar';
import type { GroupInfoResp } from '@/types';

export default function GroupSearch() {
  const navigate = useNavigate();
  const { createSession } = useChatStore();
  const [keyword, setKeyword] = useState('');
  const [results, setResults] = useState<GroupInfoResp[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    if (!keyword.trim()) return;
    setLoading(true);
    try {
      const resp = await api.searchGroupList(keyword.trim());
      if (resp.code === 0) {
        setResults(resp.data || []);
      } else {
        alert(resp.message);
      }
    } finally {
      setLoading(false);
    }
  };

  const openGroupChat = async (group: GroupInfoResp) => {
    const ok = await createSession(group.group_id || 0, 2, group.name || '');
    if (ok) {
      navigate('/chat');
    }
  };

  return (
    <div style={styles.wrap}>
      <div style={styles.card}>
        <div style={styles.header}>
          <div>
            <div style={styles.title}>搜索群聊</div>
            <div style={styles.sub}>独立页面，避免遮挡群组列表</div>
          </div>
          <button type="button" style={styles.backBtn} onClick={() => navigate('/groups')}>
            返回群组
          </button>
        </div>

        <div style={styles.searchBar}>
          <input
            style={styles.input}
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="输入群组名称关键字"
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
          <button type="button" style={styles.searchBtn} onClick={handleSearch} disabled={loading}>
            {loading ? '搜索中' : '搜索'}
          </button>
        </div>

        <div style={styles.list}>
          {results.length === 0 && <div style={styles.empty}>{keyword ? '没有找到匹配群组' : '输入关键字开始搜索'}</div>}
          {results.map((group) => {
            const groupName = group.name || `群组 ${group.group_id}`;
            return (
              <div key={group.group_id} style={styles.item}>
                <AvatarView src={resolveAvatarUrl(group.avatar)} alt={groupName} size={42} radius={12} />
                <div style={styles.meta}>
                  <div style={styles.name}>{groupName}</div>
                  <div style={styles.subMeta}>{group.member_count || 0} 人</div>
                </div>
                <button type="button" style={styles.enterBtn} onClick={() => openGroupChat(group)}>
                  进入群聊
                </button>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  wrap: { flex: 1, minHeight: 0, padding: 16, overflow: 'auto', background: 'linear-gradient(180deg, rgba(245,249,255,1) 0%, rgba(238,244,255,1) 100%)' },
  card: { maxWidth: 960, margin: '0 auto', borderRadius: 24, padding: 24, background: 'rgba(255,255,255,0.92)', border: '1px solid #dbe7fb', boxShadow: '0 18px 45px rgba(31, 64, 122, 0.08)' },
  header: { display: 'flex', justifyContent: 'space-between', gap: 16, alignItems: 'center', marginBottom: 18 },
  title: { fontSize: 22, fontWeight: 800, color: '#1f2d3d' },
  sub: { marginTop: 6, fontSize: 13, color: '#7a869a' },
  backBtn: { height: 36, padding: '0 14px', borderRadius: 10, border: '1px solid #dbe7fb', background: '#f7fbff', color: '#4a90d9', fontWeight: 700, cursor: 'pointer' },
  searchBar: { display: 'flex', gap: 10, marginBottom: 18 },
  input: { flex: 1, borderRadius: 14, border: '1px solid #dbe7fb', background: '#f8fbff', padding: '12px 14px', fontSize: 14, color: '#1f2d3d' },
  searchBtn: { height: 42, padding: '0 18px', borderRadius: 12, border: 'none', background: '#4a90d9', color: '#fff', fontWeight: 800, cursor: 'pointer' },
  list: { display: 'grid', gap: 10 },
  empty: { padding: 30, textAlign: 'center', color: '#7a869a', fontSize: 14 },
  item: { display: 'flex', alignItems: 'center', gap: 12, padding: '12px 14px', borderRadius: 16, border: '1px solid #dbe7fb', background: '#fff' },
  meta: { flex: 1, minWidth: 0 },
  name: { fontSize: 15, fontWeight: 700, color: '#1f2d3d' },
  subMeta: { marginTop: 4, fontSize: 12, color: '#7a869a' },
  enterBtn: { height: 34, padding: '0 14px', borderRadius: 10, border: '1px solid #dbe7fb', background: '#eef5ff', color: '#4a90d9', fontWeight: 700, cursor: 'pointer' },
};
