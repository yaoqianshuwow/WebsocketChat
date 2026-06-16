import { useEffect, useMemo, useState } from 'react';
import { useChatStore } from '@/store/chat';
import { useMobile } from '@/hooks/useMobile';
import AvatarView from '@/components/AvatarView';
import { resolveAvatarUrl } from '@/utils/avatar';
import { PPT_ASSISTANT_NAME, PPT_ASSISTANT_ROUTE } from '@/utils/pptAssistant';
import api from '@/api/client';

type FilterKey = 'all' | 'single' | 'group';

function safeLabel(text: string | undefined, fallback: string) {
  if (!text) return fallback;
  if (text.includes('锟?')) return fallback;
  return text;
}

export default function SessionList() {
  const isMobile = useMobile();
  const { sessions, currentSession, setCurrentSession, loadSessions, loadContacts, contacts, wsState, wsHint, peerTyping } = useChatStore();
  const [filter, setFilter] = useState<FilterKey>('all');
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  useEffect(() => {
    loadSessions();
    loadContacts();
  }, [loadSessions, loadContacts]);

  const contactMap = useMemo(() => new Map(contacts.map((item) => [item.contactId, item])), [contacts]);

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

  const openAssistant = () => {
    window.location.assign(PPT_ASSISTANT_ROUTE);
  };

  const handleDeleteSession = async (sessionId: number) => {
    const resp = await api.deleteSession(sessionId);
    if (resp.code === 0) {
      if (currentSession?.sessionId === sessionId) {
        setCurrentSession(null);
      }
      await loadSessions();
    } else {
      alert(resp.message);
    }
    setConfirmDelete(null);
  };

  return (
    <div style={{ ...styles.container, ...(isMobile ? styles.containerMobile : {}) }}>
      <div style={styles.header}>
        <div style={styles.title}>会话列表</div>
        <button type="button" style={styles.headerSearchBtn} onClick={() => window.location.assign('/messages/search')}>搜消息</button>
      </div>

      <button type="button" style={styles.assistantCard} onClick={openAssistant}>
        <div style={styles.avatarWrap}>
          <div style={styles.assistantAvatar}>
            <span style={styles.assistantAvatarText}>AI</span>
          </div>
          <span style={styles.assistantBadge}>AI</span>
        </div>
        <div style={styles.assistantInfo}>
          <div style={styles.assistantTitle}>{PPT_ASSISTANT_NAME}</div>
          <div style={styles.assistantSub}>打开 codingagent 工作台，生成和预览 PPT</div>
        </div>
        <div style={styles.assistantArrow}>&gt;</div>
      </button>

      {isMobile ? (
        <div style={styles.mobileStatusBar}>
          <span style={{ ...styles.statusDot, ...(wsState === 'connected' ? styles.statusDotOnline : wsState === 'reconnecting' || wsState === 'connecting' ? styles.statusDotBusy : styles.statusDotOffline) }} />
          <span style={styles.mobileStatusText}>{wsHint}</span>
          <span style={styles.mobileStatusMeta}>{wsState.toUpperCase()}</span>
        </div>
      ) : (
        <div style={styles.tabs}>
          <button type="button" style={{ ...styles.tab, ...(filter === 'all' ? styles.tabActive : {}) }} onClick={() => setFilter('all')}>全部</button>
          <button type="button" style={{ ...styles.tab, ...(filter === 'single' ? styles.tabActive : {}) }} onClick={() => setFilter('single')}>单聊</button>
          <button type="button" style={{ ...styles.tab, ...(filter === 'group' ? styles.tabActive : {}) }} onClick={() => setFilter('group')}>群聊</button>
        </div>
      )}

      <div style={{ ...styles.list, ...(isMobile ? styles.listMobile : {}) }}>
        {filtered.map((session) => {
          const fallbackName = session.sessionType === 2 ? `群聊 ${session.sessionId}` : `会话 ${session.sessionId}`;
          const contact = session.sessionType === 1 ? contactMap.get(session.peerId) : null;
          const displayName = safeLabel(contact?.nickname || session.sessionName, fallbackName);
          const displayAvatar = resolveAvatarUrl(contact?.avatar || session.avatar);
          const previewText = currentSession?.sessionId === session.sessionId && peerTyping ? '对方正在输入...' : safeLabel(session.lastMsgContent, '暂无消息');

          return (
            <div key={session.sessionId} style={{ display: 'flex', alignItems: 'center', borderBottom: '1px solid #f3f5f7' }}>
              <button type="button" style={{ ...styles.item, ...(isMobile ? styles.itemMobile : {}), ...(currentSession?.sessionId === session.sessionId ? styles.itemActive : {}) }} onClick={() => setCurrentSession(session)}>
                <div style={styles.avatarWrap}>
                  <AvatarView src={displayAvatar} alt={displayName} size={44} radius={12} />
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
              <button type="button" style={styles.deleteBtn} onClick={(e) => { e.stopPropagation(); setConfirmDelete(session.sessionId); }}>删除</button>
            </div>
          );
        })}
        {filtered.length === 0 && <div style={styles.empty}>暂无会话</div>}
      </div>

      {confirmDelete !== null && (
        <div style={styles.modal} onClick={() => setConfirmDelete(null)}>
          <div style={styles.confirmBox} onClick={(e) => e.stopPropagation()}>
            <h4 style={styles.confirmTitle}>删除会话</h4>
            <p style={styles.confirmMessage}>确定要删除该会话吗？删除后不可恢复。</p>
            <div style={styles.confirmActions}>
              <button type="button" style={styles.confirmOk} onClick={() => handleDeleteSession(confirmDelete)}>确定</button>
              <button type="button" style={styles.confirmCancel} onClick={() => setConfirmDelete(null)}>取消</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { width: 256, borderRight: '1px solid #dbe7fb', background: 'rgba(255,255,255,0.94)', display: 'flex', flexDirection: 'column' },
  containerMobile: { width: '100%', borderRight: 'none', borderRadius: 18, overflow: 'hidden', border: '1px solid #dbe7fb', maxHeight: '100%' },
  header: { padding: '10px 12px', borderBottom: '1px solid #dbe7fb', display: 'flex', alignItems: 'center', justifyContent: 'flex-start' },
  title: { fontWeight: 700, fontSize: 14, color: '#333', flex: 1 },
  headerSearchBtn: { height: 28, padding: '0 10px', borderRadius: 8, border: '1px solid #dbe7fb', background: '#f7fbff', color: '#4a90d9', fontSize: 12, fontWeight: 600, cursor: 'pointer' },
  assistantCard: { margin: '10px 10px 8px', borderRadius: 16, border: '1px solid #dbe7fb', background: 'linear-gradient(135deg, #f7fbff 0%, #eef5ff 100%)', display: 'flex', alignItems: 'center', gap: 10, padding: '10px 12px', cursor: 'pointer', textAlign: 'left', boxShadow: '0 10px 24px rgba(31, 64, 122, 0.06)' },
  assistantInfo: { flex: 1, minWidth: 0 },
  assistantTitle: { fontSize: 14, fontWeight: 700, color: '#1f2d3d' },
  assistantSub: { marginTop: 3, fontSize: 12, color: '#6f7f95', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' },
  assistantArrow: { fontSize: 18, color: '#4a90d9', fontWeight: 700, flexShrink: 0 },
  tabs: { display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 6, padding: '8px 10px', borderBottom: '1px solid #edf3ff' },
  mobileStatusBar: { display: 'flex', alignItems: 'center', gap: 8, padding: '10px 12px', borderBottom: '1px solid #edf3ff', background: '#f8fbff' },
  mobileStatusText: { fontSize: 12, fontWeight: 700, color: '#1f2d3d' },
  mobileStatusMeta: { marginLeft: 'auto', fontSize: 11, color: '#7a869a' },
  statusDot: { width: 10, height: 10, borderRadius: 999, display: 'inline-block', flexShrink: 0 },
  statusDotOnline: { background: '#13c26b', boxShadow: '0 0 0 4px rgba(19,194,107,0.16)' },
  statusDotBusy: { background: '#ffb020', boxShadow: '0 0 0 4px rgba(255,176,32,0.16)' },
  statusDotOffline: { background: '#a0aec0', boxShadow: '0 0 0 4px rgba(160,174,192,0.14)' },
  tab: { height: 30, borderRadius: 999, border: '1px solid #d9e8fb', background: '#f7fbff', color: '#4a90d9', cursor: 'pointer', fontSize: 12 },
  tabActive: { background: '#4a90d9', color: '#fff', borderColor: '#4a90d9' },
  list: { flex: 1, overflow: 'auto', background: '#fff' },
  listMobile: { display: 'block', overflowY: 'auto', overflowX: 'hidden', flex: 1, background: '#fff' },
  deleteBtn: { height: 24, padding: '0 8px', borderRadius: 6, border: '1px solid rgba(255,77,79,0.3)', background: 'rgba(255,255,255,0.9)', color: '#ff4d4f', fontSize: 12, cursor: 'pointer', flexShrink: 0, lineHeight: '22px', marginRight: 8 },
  item: { flex: 1, minWidth: 0, display: 'flex', gap: 10, padding: '12px 8px 11px', border: 'none', borderBottom: '1px solid #f3f5f7', background: '#fff', cursor: 'pointer', textAlign: 'left', alignItems: 'center' },
  itemMobile: { minWidth: 0, width: '100%', borderRight: 'none', borderBottom: '1px solid #f3f5f7' },
  itemActive: { background: '#f4f8ff' },
  avatarWrap: { position: 'relative', width: 44, minWidth: 44, height: 44, borderRadius: 12, overflow: 'hidden', background: '#4a90d9', flexShrink: 0 },
  avatarImg: { width: '100%', height: '100%', objectFit: 'cover', display: 'block' },
  assistantBadge: { position: 'absolute', right: -4, top: -4, minWidth: 18, height: 18, padding: '0 5px', borderRadius: 999, background: '#ffb020', color: '#fff', fontSize: 10, fontWeight: 700, display: 'grid', placeItems: 'center', border: '2px solid #fff' },
  avatarBadge: { position: 'absolute', right: -4, top: -4, minWidth: 18, height: 18, padding: '0 5px', borderRadius: 999, background: '#ff3b30', color: '#fff', fontSize: 10, fontWeight: 700, display: 'grid', placeItems: 'center', border: '2px solid #fff' },
  info: { flex: 1, minWidth: 0 },
  top: { display: 'flex', justifyContent: 'space-between', gap: 10, alignItems: 'center' },
  name: { fontSize: 15, fontWeight: 600, color: '#1f2d3d', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' },
  time: { fontSize: 11, color: '#a0a7b4', flexShrink: 0 },
  bottom: { display: 'flex', justifyContent: 'space-between', gap: 8, marginTop: 4, alignItems: 'center' },
  preview: { fontSize: 13, color: '#8a94a6', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', flex: 1 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
  modal: { position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 },
  confirmBox: { background: '#fff', borderRadius: 12, width: 300, maxWidth: '90vw', padding: 24 },
  confirmTitle: { margin: '0 0 12px', fontSize: 16, fontWeight: 600, color: '#333' },
  confirmMessage: { fontSize: 14, color: '#666', marginBottom: 20, lineHeight: 1.5 },
  confirmActions: { display: 'flex', gap: 8, justifyContent: 'flex-end' },
  confirmOk: { padding: '8px 20px', background: '#ff4d4f', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  confirmCancel: { padding: '8px 20px', background: '#f5f5f5', color: '#666', border: '1px solid #e0e0e0', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
};
