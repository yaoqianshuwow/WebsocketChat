import { useEffect, useMemo, useState } from 'react';
import { useChatStore } from '@/store/chat';
import { useMobile } from '@/hooks/useMobile';
import { pickAvatar } from '@/utils/avatar';

type FilterKey = 'all' | 'single' | 'group';

function safeLabel(text: string | undefined, fallback: string) {
  if (!text) return fallback;
  if (text.includes('�')) return fallback;
  return text;
}

export default function SessionList() {
  const isMobile = useMobile();
  const { sessions, currentSession, setCurrentSession, loadSessions, loadContacts, contacts, wsState, wsHint } = useChatStore();
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
    const now = new Date();
    const sameDay = now.toDateString() === d.toDateString();
    if (sameDay) {
      return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
    }
    return `${String(d.getMonth() + 1).padStart(2, '0')}/${String(d.getDate()).padStart(2, '0')}`;
  };

  return (
    <div style={{ ...styles.container, ...(isMobile ? styles.containerMobile : {}) }}>
      <div style={styles.header}>
        <div style={styles.title}>会话列表</div>
      </div>

      {isMobile ? (
        <div style={styles.mobileStatusBar}>
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
          <span style={styles.mobileStatusText}>{wsHint}</span>
          <span style={styles.mobileStatusMeta}>{wsState.toUpperCase()}</span>
        </div>
      ) : (
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
      )}

      <div style={{ ...styles.list, ...(isMobile ? styles.listMobile : {}) }}>
        {filtered.map((session) => {
          const fallbackName = session.sessionType === 2 ? `群聊 ${session.sessionId}` : `会话 ${session.sessionId}`;
          const contact = session.sessionType === 1 ? contactMap.get(session.peerId) : null;
          const displayName = safeLabel(contact?.nickname || session.sessionName, fallbackName);
          const displayAvatar = pickAvatar(session.peerId, contact?.avatar || session.avatar);
          const previewText = safeLabel(session.lastMsgContent, '暂无消息');
          return (
            <button
              key={session.sessionId}
              type="button"
              style={{
                ...styles.item,
                ...(isMobile ? styles.itemMobile : {}),
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
                {session.unreadCount > 0 && <span style={styles.avatarBadge}>{session.unreadCount > 99 ? '99+' : session.unreadCount}</span>}
              </div>
              <div style={styles.info}>
                <div style={styles.top}>
                  <span style={styles.name}>{displayName}</span>
                  <span style={styles.time}>{formatTime(session.lastMsgTime)}</span>
                </div>
                <div style={styles.bottom}>
                  <span style={styles.preview}>{previewText}</span>
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
    width: 256,
    borderRight: '1px solid #dbe7fb',
    background: 'rgba(255,255,255,0.94)',
    display: 'flex',
    flexDirection: 'column',
  },
  containerMobile: {
    width: '100%',
    borderRight: 'none',
    borderRadius: 18,
    overflow: 'hidden',
    border: '1px solid #dbe7fb',
    maxHeight: '100%',
  },
  header: {
    padding: '10px 12px',
    borderBottom: '1px solid #dbe7fb',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-start',
  },
  title: { fontWeight: 700, fontSize: 14, color: '#333' },
  tabs: {
    display: 'grid',
    gridTemplateColumns: 'repeat(3, 1fr)',
    gap: 6,
    padding: '8px 10px',
    borderBottom: '1px solid #edf3ff',
  },
  mobileStatusBar: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    padding: '10px 12px',
    borderBottom: '1px solid #edf3ff',
    background: '#f8fbff',
  },
  mobileStatusText: { fontSize: 12, fontWeight: 700, color: '#1f2d3d' },
  mobileStatusMeta: { marginLeft: 'auto', fontSize: 11, color: '#7a869a' },
  statusDot: { width: 10, height: 10, borderRadius: 999, display: 'inline-block', flexShrink: 0 },
  statusDotOnline: { background: '#13c26b', boxShadow: '0 0 0 4px rgba(19,194,107,0.16)' },
  statusDotBusy: { background: '#ffb020', boxShadow: '0 0 0 4px rgba(255,176,32,0.16)' },
  statusDotOffline: { background: '#a0aec0', boxShadow: '0 0 0 4px rgba(160,174,192,0.14)' },
  tab: {
    height: 30,
    borderRadius: 999,
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
  list: { flex: 1, overflow: 'auto', background: '#fff' },
  listMobile: { display: 'block', overflowY: 'auto', overflowX: 'hidden', flex: 1, background: '#fff' },
  item: {
    width: '100%',
    display: 'flex',
    gap: 12,
    padding: '12px 12px 11px',
    border: 'none',
    borderBottom: '1px solid #f3f5f7',
    background: '#fff',
    cursor: 'pointer',
    textAlign: 'left',
    alignItems: 'center',
  },
  itemMobile: {
    minWidth: 0,
    width: '100%',
    borderRight: 'none',
    borderBottom: '1px solid #f3f5f7',
  },
  itemActive: { background: '#f4f8ff' },
  avatarWrap: {
    position: 'relative',
    width: 44,
    height: 44,
    borderRadius: 12,
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
  avatarBadge: {
    position: 'absolute',
    right: -4,
    top: -4,
    minWidth: 18,
    height: 18,
    padding: '0 5px',
    borderRadius: 999,
    background: '#ff3b30',
    color: '#fff',
    fontSize: 10,
    fontWeight: 700,
    display: 'grid',
    placeItems: 'center',
    border: '2px solid #fff',
  },
  info: { flex: 1, minWidth: 0 },
  top: { display: 'flex', justifyContent: 'space-between', gap: 10, alignItems: 'center' },
  name: { fontSize: 15, fontWeight: 600, color: '#1f2d3d', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' },
  time: { fontSize: 11, color: '#a0a7b4', flexShrink: 0 },
  bottom: { display: 'flex', justifyContent: 'space-between', gap: 8, marginTop: 4, alignItems: 'center' },
  preview: {
    fontSize: 13,
    color: '#8a94a6',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    flex: 1,
  },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
};
