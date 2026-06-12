import { useState, useCallback } from 'react';
import { useChatStore } from '@/store/chat';
import wsClient from '@/ws/client';

export default function MessageInput() {
  const [text, setText] = useState('');
  const { currentSession } = useChatStore();

  const handleSend = useCallback(() => {
    const content = text.trim();
    if (!content || !currentSession) return;

    wsClient.sendMessage(content, currentSession.peerId, currentSession.sessionType);
    setText('');
  }, [text, currentSession]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  if (!currentSession) {
    return null;
  }

  return (
    <div style={styles.container}>
      <textarea
        style={styles.input}
        value={text}
        onChange={(e) => setText(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="输入消息..."
        rows={3}
      />
      <button
        style={styles.sendBtn}
        onClick={handleSend}
        disabled={!text.trim()}
      >
        发送
      </button>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: 'flex', alignItems: 'flex-end', gap: 8, padding: '8px 16px',
    borderTop: '1px solid #e8e8e8', background: '#fff',
  },
  input: {
    flex: 1, resize: 'none', padding: '8px 12px', borderRadius: 6, border: '1px solid #e0e0e0',
    fontSize: 14, lineHeight: 1.5, outline: 'none',
  },
  sendBtn: {
    padding: '8px 20px', background: '#4a90d9', color: '#fff', border: 'none',
    borderRadius: 6, cursor: 'pointer', fontSize: 14, fontWeight: 500,
    height: 40,
  },
};
