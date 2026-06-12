import { useEffect, useState } from 'react';
import api from '@/api/client';
import { useChatStore } from '@/store/chat';
import type { ContactVo, ApplyVo } from '@/types';

export default function Contacts() {
  const { contacts, loadContacts } = useChatStore();
  const [applyList, setApplyList] = useState<ApplyVo[]>([]);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [tab, setTab] = useState<'list' | 'apply' | 'search'>('list');

  useEffect(() => {
    loadContacts();
    loadApplies();
  }, []);

  const loadApplies = async () => {
    const resp = await api.getApplyList();
    if (resp.code === 0) {
      setApplyList(resp.data || []);
    }
  };

  const handleApply = async (toId: number) => {
    const resp = await api.applyContact(toId, '你好，加个好友');
    if (resp.code === 0) {
      alert('好友申请已发送');
    } else {
      alert(resp.message);
    }
  };

  const handlePass = async (applyId: number, status: number) => {
    await api.passContactApply(applyId, status);
    loadApplies();
    loadContacts();
  };

  const handleSearch = async () => {
    if (!searchKeyword.trim()) return;
    const resp = await api.searchUsers(searchKeyword);
    if (resp.code === 0) {
      setSearchResults(resp.data || []);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.tabs}>
        <button style={{...styles.tab, ...(tab === 'list' ? styles.tabActive : {})}} onClick={() => setTab('list')}>
          联系人 {contacts.length > 0 && `(${contacts.length})`}
        </button>
        <button style={{...styles.tab, ...(tab === 'apply' ? styles.tabActive : {})}} onClick={() => setTab('apply')}>
          新朋友 {applyList.length > 0 && `(${applyList.length})`}
        </button>
        <button style={{...styles.tab, ...(tab === 'search' ? styles.tabActive : {})}} onClick={() => setTab('search')}>
          查找用户
        </button>
      </div>

      <div style={styles.content}>
        {tab === 'list' && (
          <div>
            {contacts.map((c: ContactVo) => (
              <div key={c.contactId} style={styles.item}>
                <div style={styles.avatar}>{c.nickname?.charAt(0) || '?'}</div>
                <div style={styles.itemInfo}>
                  <div style={styles.itemName}>{c.nickname || `用户${c.contactId}`}</div>
                  <div style={styles.itemStatus}>{c.status === 0 ? '在线' : '离线'}</div>
                </div>
              </div>
            ))}
            {contacts.length === 0 && <div style={styles.empty}>暂无联系人</div>}
          </div>
        )}

        {tab === 'apply' && (
          <div>
            {applyList.map((a: ApplyVo) => (
              <div key={a.applyId} style={styles.item}>
                <div style={styles.avatar}>{a.nickname?.charAt(0) || '?'}</div>
                <div style={styles.itemInfo}>
                  <div style={styles.itemName}>{a.nickname || `用户${a.fromId}`}</div>
                  <div style={styles.itemStatus}>{a.remark || '请求加为好友'}</div>
                </div>
                {a.status === 0 ? (
                  <div style={styles.actions}>
                    <button style={styles.passBtn} onClick={() => handlePass(a.applyId, 1)}>同意</button>
                    <button style={styles.rejectBtn} onClick={() => handlePass(a.applyId, 2)}>拒绝</button>
                  </div>
                ) : (
                  <span style={styles.doneText}>{a.status === 1 ? '已同意' : '已拒绝'}</span>
                )}
              </div>
            ))}
            {applyList.length === 0 && <div style={styles.empty}>暂无新朋友申请</div>}
          </div>
        )}

        {tab === 'search' && (
          <div>
            <div style={styles.searchBar}>
              <input
                style={styles.searchInput}
                value={searchKeyword}
                onChange={(e) => setSearchKeyword(e.target.value)}
                placeholder="搜索用户名/昵称"
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              />
              <button style={styles.searchBtn} onClick={handleSearch}>搜索</button>
            </div>
            {searchResults.map((u: any) => (
              <div key={u.user_id} style={styles.item}>
                <div style={styles.avatar}>{u.nickname?.charAt(0) || u.username?.charAt(0) || '?'}</div>
                <div style={styles.itemInfo}>
                  <div style={styles.itemName}>{u.nickname || u.username}</div>
                  <div style={styles.itemStatus}>{u.bio || ''}</div>
                </div>
                <button style={styles.addBtn} onClick={() => handleApply(u.user_id)}>加好友</button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { flex: 1, display: 'flex', flexDirection: 'column', background: '#fff' },
  tabs: { display: 'flex', borderBottom: '1px solid #e8e8e8' },
  tab: { flex: 1, padding: '12px 0', border: 'none', background: '#fafafa', cursor: 'pointer', fontSize: 14, color: '#666', borderBottom: '2px solid transparent' },
  tabActive: { background: '#fff', color: '#4a90d9', fontWeight: 600, borderBottomColor: '#4a90d9' },
  content: { flex: 1, overflow: 'auto', padding: 8 },
  item: { display: 'flex', alignItems: 'center', padding: '10px 12px', borderBottom: '1px solid #f0f0f0' },
  avatar: { width: 36, height: 36, borderRadius: 8, background: '#4a90d9', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 600, flexShrink: 0 },
  itemInfo: { flex: 1, marginLeft: 10 },
  itemName: { fontSize: 14, fontWeight: 500, color: '#333' },
  itemStatus: { fontSize: 12, color: '#999', marginTop: 2 },
  actions: { display: 'flex', gap: 6 },
  passBtn: { padding: '4px 12px', background: '#52c41a', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  rejectBtn: { padding: '4px 12px', background: '#ff4d4f', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  doneText: { fontSize: 12, color: '#999' },
  addBtn: { padding: '4px 12px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  searchBar: { display: 'flex', gap: 8, padding: '8px 12px' },
  searchInput: { flex: 1, padding: '8px 12px', borderRadius: 6, border: '1px solid #e0e0e0', fontSize: 14, outline: 'none' },
  searchBtn: { padding: '8px 16px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
};
