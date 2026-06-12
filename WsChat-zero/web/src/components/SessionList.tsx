import { useEffect } from 'react';
import { useChatStore } from '@/store/chat';

export default function SessionList() {
  const { sessions, currentSession, setCurrentSession, loadSessions } = useChatStore();

  useEffect(() => {
    loadSessions();
  }, []);

  const formatTime = (ts?: number) => {
    if (!ts) return '';
    const d = new Date(ts * 1000);
    return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`;
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>会话列表</div>
      <div style={styles.list}>
        {sessions.map((s) => (
          <div
            key={s.sessionId}
            style={{
              ...styles.item,
              ...(currentSession?.sessionId === s.sessionId ? styles.itemActive : {}),
            }}
            onClick={() => setCurrentSession(s)}
          >
            <div style={styles.avatar}>
              {s.sessionName?.charAt(0) || '?'}
            </div>
            <div style={styles.info}>
              <div style={styles.top}>
                <span style={styles.name}>{s.sessionName || `会话${s.sessionId}`}</span>
                <span style={styles.time}>{formatTime(s.lastMsgTime)}</span>
              </div>
              <div style={styles.bottom}>
                <span style={styles.preview}>{s.lastMsgContent || '暂无消息'}</span>
                {s.unreadCount > 0 && (
                  <span style={styles.badge}>{s.unreadCount}</span>
                )}
              </div>
            </div>
          </div>
        ))}
        {sessions.length === 0 && <div style={styles.empty}>暂无会话</div>}
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { width: 280, borderRight: '1px solid #e8e8e8', background: '#fff', display: 'flex', flexDirection: 'column' },
  header: { padding: '12px 16px', fontWeight: 600, fontSize: 14, borderBottom: '1px solid #e8e8e8', color: '#333' },
  list: { flex: 1, overflow: 'auto' },
  item: { display: 'flex', padding: '10px 16px', cursor: 'pointer', borderBottom: '1px solid #f0f0f0', transition: 'background 0.2s' },
  itemActive: { background: '#e6f0ff' },
  avatar: {
    width: 40, height: 40, borderRadius: 8, background: '#4a90d9', color: '#fff',
    display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 600, fontSize: 16, flexShrink: 0,
  },
  info: { flex: 1, marginLeft: 10, minWidth: 0 },
  top: { display: 'flex', justifyContent: 'space-between', alignItems: 'center' },
  name: { fontWeight: 500, fontSize: 14, color: '#333' },
  time: { fontSize: 11, color: '#999' },
  bottom: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 4 },
  preview: { fontSize: 12, color: '#999', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', flex: 1 },
  badge: { background: '#f00', color: '#fff', borderRadius: 10, padding: '1px 6px', fontSize: 11, marginLeft: 4 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
};
