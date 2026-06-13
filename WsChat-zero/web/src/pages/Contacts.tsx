import { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '@/api/client';
import { useChatStore } from '@/store/chat';
import { pickAvatar } from '@/utils/avatar';
import type { ApplyVo, ContactVo, UserInfoResp } from '@/types';

type TabKey = 'list' | 'search';
function safeText(value: string | undefined, fallback: string) {
  if (!value) return fallback;
  if (value.includes('�')) return fallback;
  return value;
}

function displayName(name?: string, fallback = '用户') {
  return safeText(name, fallback);
}

export default function Contacts() {
  const { contacts, loadContacts, createSession } = useChatStore();
  const navigate = useNavigate();
  const [applyList, setApplyList] = useState<ApplyVo[]>([]);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [searchResults, setSearchResults] = useState<UserInfoResp[]>([]);
  const [tab, setTab] = useState<TabKey>('list');
  const [applyOpen, setApplyOpen] = useState(false);
  const [quickOpen, setQuickOpen] = useState(false);

  useEffect(() => {
    void refresh();
  }, []);

  const refresh = async () => {
    await Promise.all([loadContacts(), loadApplies()]);
  };

  const loadApplies = async () => {
    const resp = await api.getApplyList();
    if (resp.code === 0) {
      setApplyList(resp.data || []);
    }
  };

  const openChat = async (contactId: number, sessionName?: string) => {
    const ok = await createSession(contactId, 1, sessionName || '');
    if (ok) {
      navigate('/chat');
    }
  };

  const handleApply = async (toId: number) => {
    const resp = await api.applyContact(toId, '你好，想加个好友');
    alert(resp.message);
    if (resp.code === 0) {
      await loadApplies();
    }
  };

  const handlePass = async (applyId: number, status: number) => {
    const resp = await api.passContactApply(applyId, status);
    alert(resp.message);
    if (resp.code === 0) {
      await refresh();
    }
  };

  const handleSearch = async () => {
    if (!searchKeyword.trim()) return;
    const resp = await api.searchUsers(searchKeyword.trim());
    if (resp.code === 0) {
      setSearchResults(resp.data || []);
      setTab('search');
      setQuickOpen(false);
    } else {
      alert(resp.message);
    }
  };

  const applyCount = useMemo(() => applyList.filter((item) => item.status === 0).length, [applyList]);

  return (
    <div style={styles.container}>
      <div style={styles.tabs}>
        <button type="button" style={{ ...styles.tab, ...(tab === 'list' ? styles.tabActive : {}) }} onClick={() => setTab('list')}>
          联系人
        </button>
        <div style={styles.plusWrap}>
          <button
            type="button"
            style={{ ...styles.plusBtn, ...(applyOpen || quickOpen ? styles.plusBtnActive : {}) }}
            onClick={() => {
              setQuickOpen((v) => !v);
              setApplyOpen(false);
            }}
            title="更多操作"
          >
            +
            {applyCount > 0 && <span style={styles.badge}>{applyCount}</span>}
          </button>
          {quickOpen && (
            <div style={styles.quickPanel}>
              <button
                type="button"
                style={styles.quickItem}
                onClick={() => {
                  setApplyOpen(true);
                  setTab('list');
                  setQuickOpen(false);
                }}
              >
                <span>新朋友</span>
                {applyCount > 0 && <span style={styles.quickBadge}>{applyCount}</span>}
              </button>
              <button
                type="button"
                style={styles.quickItem}
                onClick={() => {
                  setTab('search');
                  setApplyOpen(false);
                  setQuickOpen(false);
                }}
              >
                <span>查找用户</span>
              </button>
            </div>
          )}
        </div>
      </div>

      {applyOpen && (
        <div style={styles.applyPanel}>
          <div style={styles.applyHeader}>
            <span>新朋友</span>
            <button type="button" style={styles.linkBtn} onClick={refresh}>
              刷新
            </button>
          </div>
          <div style={styles.applyList}>
            {applyList.map((item: ApplyVo) => {
              const name = displayName(item.nickname, `用户 ${item.fromId}`);
              const avatarSrc = pickAvatar(item.fromId, undefined);
              return (
                <div key={item.applyId} style={styles.item}>
                  <img
                    src={avatarSrc}
                    alt={name}
                    style={styles.avatarImg}
                    onError={(e) => {
                      (e.currentTarget as HTMLImageElement).src = pickAvatar(item.fromId + 1);
                    }}
                  />
                  <div style={styles.itemInfo}>
                    <div style={styles.itemName}>{name}</div>
                    <div style={styles.itemStatus}>{safeText(item.remark, '好友申请')}</div>
                  </div>
                  {item.status === 0 ? (
                    <div style={styles.actions}>
                      <button type="button" style={styles.passBtn} onClick={() => handlePass(item.applyId, 1)}>
                        同意
                      </button>
                      <button type="button" style={styles.rejectBtn} onClick={() => handlePass(item.applyId, 2)}>
                        拒绝
                      </button>
                    </div>
                  ) : (
                    <span style={styles.doneText}>{item.status === 1 ? '已同意' : '已拒绝'}</span>
                  )}
                </div>
              );
            })}
            {applyList.length === 0 && <div style={styles.empty}>暂无新的好友申请</div>}
          </div>
        </div>
      )}

      <div style={styles.content}>
        {tab === 'list' && (
          <div>
            {contacts.map((contact: ContactVo) => {
              const name = displayName(contact.nickname, '好友');
              const avatarSrc = pickAvatar(contact.contactId, contact.avatar);
              return (
                <button key={contact.contactId} type="button" style={styles.item} onClick={() => openChat(contact.contactId, contact.nickname)}>
                  <img
                    src={avatarSrc}
                    alt={name}
                    style={styles.avatarImg}
                    onError={(e) => {
                      (e.currentTarget as HTMLImageElement).src = pickAvatar(contact.contactId + 1);
                    }}
                  />
                  <div style={styles.itemInfo}>
                    <div style={styles.itemName}>{name}</div>
                    <div style={styles.itemStatus}>{contact.status === 0 ? '在线' : '离线'}</div>
                  </div>
                  <span style={styles.enter}>聊天</span>
                </button>
              );
            })}
            {contacts.length === 0 && <div style={styles.empty}>暂无联系人</div>}
          </div>
        )}

        {tab === 'search' && (
          <div>
            <div style={styles.searchBar}>
              <input
                style={styles.searchInput}
                value={searchKeyword}
                onChange={(e) => setSearchKeyword(e.target.value)}
                placeholder="搜索用户名或昵称"
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              />
              <button type="button" style={styles.searchBtn} onClick={handleSearch}>
                搜索
              </button>
            </div>
            {searchResults.map((user) => {
              const name = displayName(user.nickname || user.username, `用户 ${user.user_id}`);
              const avatarSrc = pickAvatar(user.user_id, user.avatar);
              return (
                <div key={user.user_id} style={styles.item}>
                  <img
                    src={avatarSrc}
                    alt={name}
                    style={styles.avatarImg}
                    onError={(e) => {
                      (e.currentTarget as HTMLImageElement).src = pickAvatar((user.user_id || 0) + 1);
                    }}
                  />
                  <div style={styles.itemInfo}>
                    <div style={styles.itemName}>{name}</div>
                    <div style={styles.itemStatus}>{safeText(user.bio, '暂无签名')}</div>
                  </div>
                  <button type="button" style={styles.addBtn} onClick={() => user.user_id && handleApply(user.user_id)}>
                    添加好友
                  </button>
                </div>
              );
            })}
            {searchResults.length === 0 && <div style={styles.empty}>请输入关键词后搜索</div>}
          </div>
        )}
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { flex: 1, display: 'flex', flexDirection: 'column', background: '#fff' },
  tabs: { display: 'grid', gridTemplateColumns: '1fr 64px', alignItems: 'stretch', borderBottom: '1px solid #e8e8e8' },
  tab: {
    padding: '12px 0',
    border: 'none',
    background: '#fafafa',
    cursor: 'pointer',
    fontSize: 14,
    color: '#666',
    borderBottom: '2px solid transparent',
  },
  tabActive: { background: '#fff', color: '#4a90d9', fontWeight: 600, borderBottomColor: '#4a90d9' },
  plusWrap: { position: 'relative' },
  plusBtn: {
    position: 'relative',
    width: '100%',
    height: '100%',
    border: 'none',
    background: '#fff',
    cursor: 'pointer',
    fontSize: 28,
    lineHeight: 1,
    color: '#4a90d9',
    borderBottom: '2px solid transparent',
  },
  plusBtnActive: { background: '#eaf3ff', borderBottomColor: '#4a90d9' },
  quickPanel: {
    position: 'absolute',
    top: 'calc(100% + 8px)',
    right: 8,
    minWidth: 150,
    padding: 8,
    borderRadius: 12,
    background: '#fff',
    border: '1px solid #dbe7fb',
    boxShadow: '0 16px 32px rgba(31, 64, 122, 0.14)',
    zIndex: 20,
  },
  quickItem: {
    width: '100%',
    minHeight: 40,
    padding: '0 12px',
    border: 'none',
    borderRadius: 8,
    background: '#fff',
    color: '#1f2d3d',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    fontSize: 13,
    textAlign: 'left',
  },
  quickBadge: {
    minWidth: 18,
    height: 18,
    padding: '0 5px',
    borderRadius: 9,
    background: '#f5222d',
    color: '#fff',
    fontSize: 11,
    display: 'grid',
    placeItems: 'center',
  },
  badge: {
    position: 'absolute',
    top: 8,
    right: 7,
    minWidth: 18,
    height: 18,
    padding: '0 5px',
    borderRadius: 9,
    background: '#f5222d',
    color: '#fff',
    fontSize: 11,
    display: 'grid',
    placeItems: 'center',
  },
  applyPanel: {
    borderBottom: '1px solid #e8e8e8',
    background: '#f7fbff',
  },
  applyHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '10px 12px',
    fontSize: 13,
    color: '#333',
    borderBottom: '1px solid #edf2f7',
  },
  linkBtn: {
    border: 'none',
    background: 'transparent',
    color: '#4a90d9',
    cursor: 'pointer',
    padding: 0,
  },
  applyList: { maxHeight: 220, overflow: 'auto' },
  content: { flex: 1, overflow: 'auto', padding: 8 },
  item: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    padding: '10px 12px',
    border: 'none',
    borderBottom: '1px solid #f0f0f0',
    background: '#fff',
    cursor: 'pointer',
    textAlign: 'left',
  },
  avatarImg: {
    width: 36,
    height: 36,
    borderRadius: 8,
    objectFit: 'cover',
    flexShrink: 0,
    background: '#e6f0ff',
  },
  itemInfo: { flex: 1, minWidth: 0 },
  itemName: { fontSize: 14, fontWeight: 500, color: '#333' },
  itemStatus: { fontSize: 12, color: '#999', marginTop: 2, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' },
  actions: { display: 'flex', gap: 6 },
  passBtn: { padding: '4px 12px', background: '#52c41a', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  rejectBtn: { padding: '4px 12px', background: '#ff4d4f', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  doneText: { fontSize: 12, color: '#999' },
  addBtn: { padding: '4px 12px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  enter: { fontSize: 12, color: '#4a90d9', flexShrink: 0 },
  searchBar: { display: 'flex', gap: 8, padding: '8px 12px' },
  searchInput: { flex: 1, padding: '8px 12px', borderRadius: 6, border: '1px solid #e0e0e0', fontSize: 14, outline: 'none' },
  searchBtn: { padding: '8px 16px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
};
