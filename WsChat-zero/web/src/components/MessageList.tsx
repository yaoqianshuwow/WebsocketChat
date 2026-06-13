import { useEffect, useRef } from 'react';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import type { MessageVo } from '@/types';

function safeLabel(text: string | undefined, fallback: string) {
  if (!text) return fallback;
  if (text.includes('�')) return fallback;
  return text;
}

export default function MessageList() {
  const { messages, currentSession } = useChatStore();
  const { user, userInfo } = useAuthStore();
  const bottomRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  if (!currentSession) {
    return <div style={styles.empty}>请选择一个会话开始聊天</div>;
  }

  const formatTime = (ts: number) => {
    const d = new Date(ts * 1000);
    return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
  };

  const headerTitle = safeLabel(currentSession.sessionName, currentSession.sessionType === 2 ? `群聊 ${currentSession.sessionId}` : `会话 ${currentSession.sessionId}`);
  const peerName = safeLabel(currentSession.sessionName, '对方');
  const currentUserId = user?.userId || user?.user_id || userInfo?.user_id;
  const currentNickname = safeLabel(userInfo?.nickname || user?.nickname, '我');

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <div>
          <div style={styles.headerTitle}>{headerTitle}</div>
          <div style={styles.headerSub}>消息窗口固定高度，向上滚动查看历史消息</div>
        </div>
      </div>
      <div ref={listRef} style={styles.list}>
        {messages.map((msg: MessageVo) => {
          const isMe = msg.mine || msg.senderId === currentUserId;
          const senderName = isMe ? currentNickname : peerName;
          const bubbleStyle = isMe ? styles.bubbleMe : styles.bubbleOther;
          const rowStyle = isMe ? styles.msgRowMe : styles.msgRowOther;
          return (
            <div key={msg.msgId || msg.localId} style={{ ...styles.msgRow, ...rowStyle }}>
              <div style={{ ...styles.avatar, ...(isMe ? styles.avatarMe : styles.avatarOther) }}>
                {(senderName || '?').charAt(0)}
              </div>
              <div style={{ ...styles.bubble, ...bubbleStyle }}>
                <div style={styles.sender}>{senderName}</div>
                <div style={{ ...styles.content, ...(isMe ? styles.contentMe : {}) }}>{msg.content || msg.fileName || '[非文本消息]'}</div>
                <div style={styles.metaRow}>
                  {msg.status && isMe && <span style={styles.status}>{msg.status === 'sending' ? '发送中' : msg.status === 'failed' ? '发送失败' : '已发送'}</span>}
                  <div style={styles.time}>{formatTime(msg.createdAt)}</div>
                </div>
              </div>
            </div>
          );
        })}
        <div ref={bottomRef} />
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { flex: 1, minHeight: 0, display: 'flex', flexDirection: 'column', background: '#f3f7ff' },
  header: {
    padding: '14px 18px',
    borderBottom: '1px solid #dbe7fb',
    background: 'rgba(255,255,255,0.92)',
    color: '#1f2d3d',
  },
  headerTitle: { fontWeight: 700, fontSize: 15 },
  headerSub: { marginTop: 4, fontSize: 12, color: '#7a869a' },
  list: { flex: 1, overflowY: 'auto', padding: '18px 20px', display: 'flex', flexDirection: 'column', gap: 12 },
  empty: { flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#999', fontSize: 16 },
  msgRow: { display: 'flex', gap: 10, alignItems: 'flex-end' },
  msgRowMe: { justifyContent: 'flex-end' },
  msgRowOther: { justifyContent: 'flex-start' },
  avatar: {
    width: 34,
    height: 34,
    borderRadius: 12,
    display: 'grid',
    placeItems: 'center',
    fontSize: 12,
    fontWeight: 700,
    flexShrink: 0,
  },
  avatarMe: { background: '#1f6feb', color: '#fff', order: 2 },
  avatarOther: { background: '#e6eefc', color: '#36527a' },
  bubble: { maxWidth: '62%', padding: '10px 12px', borderRadius: 18, boxShadow: '0 10px 20px rgba(31, 45, 61, 0.06)' },
  bubbleMe: { background: '#1f6feb', color: '#fff', borderBottomRightRadius: 6, order: 1 },
  bubbleOther: { background: '#fff', color: '#1f2d3d', border: '1px solid #dbe7fb', borderBottomLeftRadius: 6 },
  sender: { fontSize: 11, marginBottom: 4, opacity: 0.78 },
  content: { fontSize: 14, color: '#333', lineHeight: 1.5, wordBreak: 'break-word' },
  contentMe: { color: '#fff' },
  metaRow: { marginTop: 6, display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 10 },
  status: { fontSize: 11, opacity: 0.82 },
  time: { fontSize: 10, opacity: 0.72, textAlign: 'right' },
};
