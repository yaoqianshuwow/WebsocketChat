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
  const { user } = useAuthStore();
  const bottomRef = useRef<HTMLDivElement>(null);

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

  return (
    <div style={styles.container}>
      <div style={styles.header}>{headerTitle}</div>
      <div style={styles.list}>
        {messages.map((msg: MessageVo) => {
          const isMe = msg.senderId === user?.userId || msg.senderId === user?.user_id;
          return (
            <div key={msg.msgId || msg.localId} style={{ ...styles.msgRow, ...(isMe ? styles.msgRowMe : {}) }}>
              <div style={styles.bubble}>
                <div style={styles.sender}>{isMe ? '我' : `用户 ${msg.senderId}`}</div>
                <div style={styles.content}>{msg.content || msg.fileName || '[非文本消息]'}</div>
                <div style={styles.time}>{formatTime(msg.createdAt)}</div>
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
  container: { flex: 1, display: 'flex', flexDirection: 'column' },
  header: {
    padding: '12px 16px',
    fontWeight: 600,
    fontSize: 14,
    borderBottom: '1px solid #e8e8e8',
    background: '#fff',
    color: '#333',
  },
  list: { flex: 1, overflow: 'auto', padding: 16, background: '#fafafa' },
  empty: { flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#999', fontSize: 16 },
  msgRow: { display: 'flex', marginBottom: 12 },
  msgRowMe: { justifyContent: 'flex-end' },
  bubble: { maxWidth: '70%', padding: '8px 12px', borderRadius: 8, background: '#fff', boxShadow: '0 1px 2px rgba(0,0,0,0.1)' },
  sender: { fontSize: 11, color: '#4a90d9', marginBottom: 2 },
  content: { fontSize: 14, color: '#333', lineHeight: 1.5, wordBreak: 'break-word' },
  time: { fontSize: 10, color: '#bbb', textAlign: 'right', marginTop: 4 },
};
