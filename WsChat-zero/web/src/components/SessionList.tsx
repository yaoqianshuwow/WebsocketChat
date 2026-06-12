import { useEffect, useMemo, useState } from 'react';
import { useChatStore } from '@/store/chat';
import { pickAvatar } from '@/utils/avatar';

type FilterKey = 'all' | 'single' | 'group';

function safeLabel(text: string | undefined, fallback: string) {
  if (!text) return fallback;
  if (text.includes('�')) return fallback;
  return text;
}

export default function SessionList() {
  const { sessions, currentSession, setCurrentSession, loadSessions, loadContacts, contacts } = useChatStore();
  const [filter, setFilter] = useState<FilterKey>('all');

  useEffect(() => {
    loadSessions();
    loadContacts();
  }, [loadSessions, loadContacts]);

  const contactMap = useMemo(() => {
    return new Map(contacts.map((item) => [item.contactId, item]));
  }, [contacts]);

  const filtered = useMemo(() => {
    if (filter === 'single') return sessions.filter((item) => item.sessionType === 1);
    if (filter === 'group') return sessions.filter((item) => item.sessionType === 2);
    return sessions;
  }, [filter, sessions]);

  const formatTime = (ts?: number) => {
    if (!ts) return '';
    const d = new Date(ts * 1000);
    return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <div style={styles.title}>会话列表</div>
        <button type="button" style={styles.refreshBtn} onClick={loadSessions}>
          刷新
        </button>
      </div>

      <div style={styles.tabs}>
        <button type="button" style={{ ...styles.tab, ...(filter === 'all' ? styles.tabActive : {}) }} onClick={() => setFilter('all')}>
          全部
        </button>
        <button type="button" style={{ ...styles.tab, ...(filter === 'single' ? styles.tabActive : {}) }} onClick={() => setFilter('single')}>
          单聊
        </button>
        <button type="button" style={{ ...styles.tab, ...(filter === 'group' ? styles.tabActive : {}) }} onClick={() => setFilter('group')}>
          群聊
        </button>
      </div>

      <div style={styles.list}>
        {filtered.map((session) => {
          const fallbackName = session.sessionType === 2 ? `群聊 ${session.sessionId}` : `会话 ${session.sessionId}`;
          const contact = session.sessionType === 1 ? contactMap.get(session.peerId) : null;
          const displayName = safeLabel(contact?.nickname || session.sessionName, fallbackName);
          const displayAvatar = pickAvatar(session.peerId, contact?.avatar || session.avatar);
          return (
            <button
              key={session.sessionId}
              type="button"
              style={{
                ...styles.item,
                ...(currentSession?.sessionId === session.sessionId ? styles.itemActive : {}),
              }}
              onClick={() => setCurrentSession(session)}
            >
              <div style={styles.avatarWrap}>
                <img
                  src={displayAvatar}
                  alt={displayName}
                  style={styles.avatarImg}
                  onError={(e) => {
                    (e.currentTarget as HTMLImageElement).src = pickAvatar(session.peerId + 1);
                  }}
                />
              </div>
              <div style={styles.info}>
                <div style={styles.top}>
                  <span style={styles.name}>{displayName}</span>
                  <span style={styles.time}>{formatTime(session.lastMsgTime)}</span>
                </div>
                <div style={styles.bottom}>
                  <span style={styles.preview}>{session.lastMsgContent || '暂无消息'}</span>
                  {session.unreadCount > 0 && <span style={styles.badge}>{session.unreadCount}</span>}
                </div>
              </div>
            </button>
          );
        })}

        {filtered.length === 0 && <div style={styles.empty}>暂无会话</div>}
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    width: 280,
    borderRight: '1px solid #e8e8e8',
    background: '#fff',
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    padding: '12px 16px',
    borderBottom: '1px solid #e8e8e8',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  title: { fontWeight: 700, fontSize: 14, color: '#333' },
  refreshBtn: {
    padding: '4px 10px',
    borderRadius: 4,
    border: '1px solid #4a90d9',
    background: '#4a90d9',
    color: '#fff',
    cursor: 'pointer',
    fontSize: 12,
  },
  tabs: {
    display: 'grid',
    gridTemplateColumns: 'repeat(3, 1fr)',
    gap: 8,
    padding: 10,
    borderBottom: '1px solid #f0f0f0',
  },
  tab: {
    height: 32,
    borderRadius: 6,
    border: '1px solid #d9e8fb',
    background: '#f7fbff',
    color: '#4a90d9',
    cursor: 'pointer',
    fontSize: 12,
  },
  tabActive: {
    background: '#4a90d9',
    color: '#fff',
    borderColor: '#4a90d9',
  },
  list: { flex: 1, overflow: 'auto' },
  item: {
    width: '100%',
    display: 'flex',
    gap: 10,
    padding: '10px 12px',
    border: 'none',
    borderBottom: '1px solid #f0f0f0',
    background: '#fff',
    cursor: 'pointer',
    textAlign: 'left',
  },
  itemActive: { background: '#eaf3ff' },
  avatar: {
    width: 40,
    height: 40,
    borderRadius: 8,
    background: '#4a90d9',
    color: '#fff',
    display: 'grid',
    placeItems: 'center',
    fontWeight: 700,
    flexShrink: 0,
  },
  avatarWrap: {
    width: 40,
    height: 40,
    borderRadius: 8,
    overflow: 'hidden',
    background: '#4a90d9',
    flexShrink: 0,
  },
  avatarImg: {
    width: '100%',
    height: '100%',
    objectFit: 'cover',
    display: 'block',
  },
  info: { flex: 1, minWidth: 0 },
  top: { display: 'flex', justifyContent: 'space-between', gap: 10 },
  name: { fontSize: 14, fontWeight: 600, color: '#333', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' },
  time: { fontSize: 11, color: '#999', flexShrink: 0 },
  bottom: { display: 'flex', justifyContent: 'space-between', gap: 8, marginTop: 4 },
  preview: {
    fontSize: 12,
    color: '#999',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    flex: 1,
  },
  badge: {
    minWidth: 18,
    height: 18,
    padding: '0 6px',
    borderRadius: 9,
    background: '#f5222d',
    color: '#fff',
    fontSize: 11,
    display: 'grid',
    placeItems: 'center',
  },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
};
