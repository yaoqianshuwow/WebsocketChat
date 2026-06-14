import { useEffect, useRef } from 'react';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import api from '@/api/client';
import type { MessageVo } from '@/types';

function safeLabel(text: string | undefined, fallback: string) {
  if (!text) return fallback;
  if (text.includes('锟')) return fallback;
  return text;
}

function formatFileSize(bytes: number | undefined): string {
  if (!bytes || bytes <= 0) return '';
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function getFileIcon(fileName: string): string {
  const ext = fileName.split('.').pop()?.toLowerCase() || '';
  const iconMap: Record<string, string> = {
    pdf: '📄', doc: '📝', docx: '📝',
    xls: '📊', xlsx: '📊',
    ppt: '📽', pptx: '📽',
    zip: '📦', rar: '📦', '7z': '📦',
    jpg: '🖼', jpeg: '🖼', png: '🖼', gif: '🖼', webp: '🖼',
    mp4: '🎬', avi: '🎬', mov: '🎬',
    mp3: '🎵', wav: '🎵',
    txt: '📃',
    exe: '⚙', dmg: '⚙',
  };
  return iconMap[ext] || '📎';
}

function AvatarCircle({ src, name, size = 34 }: { src?: string; name: string; size?: number }) {
  if (src) {
    return (
      <img
        src={src}
        alt={name}
        style={{
          width: size, height: size, borderRadius: 12,
          objectFit: 'cover', flexShrink: 0,
          background: '#e6eefc',
        }}
        onError={(e) => {
          (e.target as HTMLImageElement).style.display = 'none';
          (e.target as HTMLImageElement).nextElementSibling?.classList.remove('hidden');
        }}
      />
    );
  }
  return (
    <div style={{
      width: size, height: size, borderRadius: 12,
      display: 'grid', placeItems: 'center',
      fontSize: Math.max(11, size * 0.38), fontWeight: 700,
      background: '#e6eefc', color: '#36527a',
      flexShrink: 0,
    }}>
      {(name || '?').charAt(0)}
    </div>
  );
}

export default function MessageList() {
  const { messages, currentSession } = useChatStore();
  const { user, userInfo } = useAuthStore();
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

  const isGroupChat = currentSession.sessionType === 2;
  const headerTitle = safeLabel(
    currentSession.sessionName,
    isGroupChat ? `群聊 ${currentSession.peerId}` : `会话 ${currentSession.sessionId}`,
  );
  const peerName = safeLabel(currentSession.sessionName, '对方');
  const currentUserId = user?.userId || user?.user_id || userInfo?.user_id || 0;
  const currentAvatar = userInfo?.avatar || user?.avatar || '';
  const currentNickname = safeLabel(userInfo?.nickname || user?.nickname, '我');

  const handleDownload = (msg: MessageVo) => {
    if (msg.fileUrl) {
      api.downloadFile(msg.fileUrl);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <div>
          <div style={styles.headerTitle}>{headerTitle}</div>
          <div style={styles.headerSub}>消息列表会按会话自动刷新，群聊使用同一个共享列表</div>
        </div>
      </div>
      <div style={styles.list}>
        {messages.map((msg: MessageVo) => {
          const isMe = msg.mine || msg.senderId === currentUserId;
          const senderName = isMe ? currentNickname : msg.sendName || (isGroupChat ? `用户 ${msg.senderId}` : peerName);
          const senderAvatar = isMe ? currentAvatar : msg.sendAvatar;
          const bubbleStyle = isMe ? styles.bubbleMe : styles.bubbleOther;
          const rowStyle = isMe ? styles.msgRowMe : styles.msgRowOther;
          const isFileMsg = msg.msgType === 2 && msg.fileUrl;

          return (
            <div key={msg.msgId || msg.localId} style={{ ...styles.msgRow, ...rowStyle }}>
              <div style={isMe ? styles.avatarWrapperMe : styles.avatarWrapperOther}>
                <AvatarCircle src={senderAvatar} name={senderName} />
              </div>
              <div style={{ ...styles.bubble, ...bubbleStyle }}>
                {isGroupChat && <div style={styles.sender}>{senderName}</div>}

                {isFileMsg ? (
                  <div style={styles.fileCard} onClick={() => handleDownload(msg)}>
                    <div style={styles.fileIcon}>{getFileIcon(msg.fileName || '')}</div>
                    <div style={styles.fileInfo}>
                      <div style={{ ...styles.fileName, ...(isMe ? { color: '#fff' } : {}) }}>
                        {msg.fileName || '文件'}
                      </div>
                      <div style={{ ...styles.fileSize, ...(isMe ? { color: 'rgba(255,255,255,0.7)' } : {}) }}>
                        {formatFileSize(msg.fileSize)}
                      </div>
                    </div>
                    <div style={{ ...styles.downloadBtn, ...(isMe ? { borderColor: 'rgba(255,255,255,0.4)', color: '#fff' } : {}) }}>
                      下载
                    </div>
                  </div>
                ) : (
                  <div style={{ ...styles.content, ...(isMe ? styles.contentMe : {}) }}>
                    {msg.content || '[非文本消息]'}
                  </div>
                )}

                <div style={styles.metaRow}>
                  {msg.status && isMe && (
                    <span style={styles.status}>
                      {msg.status === 'sending' ? '发送中' : msg.status === 'failed' ? '发送失败' : '已发送'}
                    </span>
                  )}
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
  avatarWrapperMe: { order: 2 },
  avatarWrapperOther: { order: 0 },
  bubble: { maxWidth: '62%', padding: '10px 12px', borderRadius: 18, boxShadow: '0 10px 20px rgba(31, 45, 61, 0.06)' },
  bubbleMe: { background: '#1f6feb', color: '#fff', borderBottomRightRadius: 6, order: 1 },
  bubbleOther: { background: '#fff', color: '#1f2d3d', border: '1px solid #dbe7fb', borderBottomLeftRadius: 6 },
  sender: { fontSize: 11, marginBottom: 4, opacity: 0.78 },
  content: { fontSize: 14, color: '#333', lineHeight: 1.5, wordBreak: 'break-word' },
  contentMe: { color: '#fff' },
  metaRow: { marginTop: 6, display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 10 },
  status: { fontSize: 11, opacity: 0.82 },
  time: { fontSize: 10, opacity: 0.72, textAlign: 'right' },
  // 文件消息卡片
  fileCard: {
    display: 'flex', alignItems: 'center', gap: 10,
    cursor: 'pointer', padding: '4px 0',
  },
  fileIcon: { fontSize: 28, flexShrink: 0 },
  fileInfo: { flex: 1, minWidth: 0 },
  fileName: { fontSize: 13, fontWeight: 600, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', color: '#333' },
  fileSize: { fontSize: 11, marginTop: 2, color: '#7a869a' },
  downloadBtn: {
    fontSize: 12, padding: '4px 10px', borderRadius: 8,
    borderWidth: 1, borderStyle: 'solid', borderColor: '#dbe7fb',
    background: 'transparent',
    cursor: 'pointer', flexShrink: 0, color: '#4a90d9',
  },
};
