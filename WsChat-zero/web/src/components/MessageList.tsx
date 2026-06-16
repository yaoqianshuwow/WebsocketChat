import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import AvatarView from '@/components/AvatarView';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import api from '@/api/client';
import { resolveAvatarUrl } from '@/utils/avatar';
import { resolveFileUrl } from '@/utils/fileUrl';
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
    pdf: '📄', doc: '📝', docx: '📝', xls: '📊', xlsx: '📊',
    ppt: '📽', pptx: '📽', zip: '📦', rar: '📦', '7z': '📦',
    mp4: '🎬', avi: '🎬', mov: '🎬', mp3: '🎵', wav: '🎵',
    txt: '📃', exe: '⚙', dmg: '⚙',
  };
  return iconMap[ext] || '📎';
}

function isImageFile(fileName: string): boolean {
  const ext = fileName.split('.').pop()?.toLowerCase() || '';
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp'].includes(ext);
}

function isPdfFile(fileName: string): boolean {
  return fileName.toLowerCase().endsWith('.pdf');
}

// 静音获取用户头像信息并缓存
const avatarCache = new Map<number, string>();
const nameCache = new Map<number, string>();

async function resolveSenderAvatar(userId: number): Promise<string> {
  if (avatarCache.has(userId)) return avatarCache.get(userId) || '';
  try {
    const resp = await api.getUserInfo(userId);
    if (resp.code === 0) {
      const url = resolveAvatarUrl(resp.avatar);
      avatarCache.set(userId, url);
      nameCache.set(userId, resp.nickname || '');
      return url;
    }
  } catch { /* ignore */ }
  avatarCache.set(userId, '');
  return '';
}

function getSenderAvatar(senderId: number): string {
  return avatarCache.get(senderId) || '';
}

async function resolveSenderName(userId: number, fallback: string): Promise<string> {
  if (nameCache.has(userId)) return nameCache.get(userId) || fallback;
  try {
    const resp = await api.getUserInfo(userId);
    if (resp.code === 0 && resp.nickname) {
      nameCache.set(userId, resp.nickname);
      return resp.nickname;
    }
  } catch { /* ignore */ }
  return fallback;
}

// 表情文本映射
const EMOJI_TEXT_MAP: [string, string][] = [
  [':)', '😀'], ['(:', '😀'], [':(', '😢'], [":'(", '😭'],
  [':D', '😁'], [';)', '😉'], [':P', '😛'], ['B)', '😎'],
  ['<3', '❤️'], [":'-(", '🥹'], ['XD', '😆'], [':|', '😐'],
];

function replaceEmojiShortcuts(text: string): string {
  let result = text;
  for (const [key, emoji] of EMOJI_TEXT_MAP) {
    result = result.split(key).join(emoji);
  }
  return result;
}

export default function MessageList() {
  const navigate = useNavigate();
  const { messages, currentSession, loadingHistory, loadOlderMessages, peerTyping } = useChatStore();
  const { user, userInfo } = useAuthStore();
  const bottomRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const stickToBottomRef = useRef(true);
  const scrollTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [showScrollbar, setShowScrollbar] = useState(false);
  // 缓存其他用户头像：senderId -> avatarUrl
  const [, forceUpdate] = useState(0);

  // 默认滚动到底部 - 首次加载用 auto，后续消息用 smooth
  const scrollToBottom = (smooth = false) => {
    const el = listRef.current;
    if (!el) return;
    el.scrollTop = el.scrollHeight;
    if (smooth) {
      bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  };

  useEffect(() => {
    if (loadingHistory) return;
    const firstLoad = messages.length > 0 && messages.length <= 100;
    if (stickToBottomRef.current) {
      requestAnimationFrame(() => scrollToBottom(!firstLoad));
    }
  }, [messages, loadingHistory]);

  // 切换会话时重置滚动粘性并隐藏滚动条
  useEffect(() => {
    stickToBottomRef.current = true;
    setShowScrollbar(false);
    if (scrollTimerRef.current) {
      clearTimeout(scrollTimerRef.current);
      scrollTimerRef.current = null;
    }
    // 清空头像缓存（切换会话时不需要清空，但防止缓存膨胀）
    if (avatarCache.size > 50) {
      avatarCache.clear();
      nameCache.clear();
    }
  }, [currentSession?.sessionId]);

  // 为每一条消息中没有头像的 sender 异步获取头像
  useEffect(() => {
    if (!currentSession) return;
    const currentUserIdVal = user?.userId || user?.user_id || userInfo?.user_id || 0;
    const uniqueSenders = new Set<number>();
    messages.forEach((msg) => {
      if (!msg.mine && msg.senderId !== currentUserIdVal && !msg.sendAvatar && !avatarCache.has(msg.senderId)) {
        uniqueSenders.add(msg.senderId);
      }
    });
    if (uniqueSenders.size === 0) return;
    let updated = false;
    Promise.allSettled([...uniqueSenders].map(async (sid) => {
      await resolveSenderAvatar(sid);
      updated = true;
    })).then(() => {
      if (updated) forceUpdate((n) => n + 1);
    });
  }, [messages, currentSession, user, userInfo]);

  const handleScrollTop = async () => {
    const el = listRef.current;
    if (!el || el.scrollTop > 60 || loadingHistory) return;

    const prevHeight = el.scrollHeight;
    const prevTop = el.scrollTop;
    await loadOlderMessages();

    requestAnimationFrame(() => {
      const next = listRef.current;
      if (!next) return;
      next.scrollTop = prevTop + (next.scrollHeight - prevHeight);
    });
  };

  const handleScrollState = () => {
    const el = listRef.current;
    if (!el) return;
    const distanceToBottom = el.scrollHeight - el.scrollTop - el.clientHeight;
    stickToBottomRef.current = distanceToBottom < 60;
    setShowScrollbar(true);
    if (scrollTimerRef.current) clearTimeout(scrollTimerRef.current);
    scrollTimerRef.current = setTimeout(() => {
      setShowScrollbar(false);
      scrollTimerRef.current = null;
    }, 1200);
  };

  if (!currentSession) {
    return <div style={styles.empty}>请选择一个会话开始聊天</div>;
  }

  const formatTime = (ts: number) => {
    const d = new Date(ts * 1000);
    return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
  };

  const isGroupChat = currentSession.sessionType === 2;
  const headerTitle = safeLabel(currentSession.sessionName, isGroupChat ? `群聊 ${currentSession.peerId}` : `会话 ${currentSession.sessionId}`);
  const currentUserIdVal = user?.userId || user?.user_id || userInfo?.user_id || 0;
  const currentAvatar = resolveAvatarUrl(userInfo?.avatar || user?.avatar || '');
  const currentNickname = safeLabel(userInfo?.nickname || user?.nickname, '我');

  const handleDownload = (msg: MessageVo) => {
    if (msg.fileUrl) api.downloadFile(msg.fileUrl);
  };

  // PDF文件URL处理：返回预览图或直接下载
  const getPdfPreviewUrl = (fileUrl?: string): string | undefined => {
    if (!fileUrl) return undefined;
    // 直接返回文件URL，浏览器支持PDF预览
    return resolveFileUrl(fileUrl);
  };

  return (
    <div style={styles.container}>
      <style>{`
        .message-scroll { scrollbar-gutter: stable; }
        .message-scroll-hidden { scrollbar-width: none; -ms-overflow-style: none; }
        .message-scroll-hidden::-webkit-scrollbar { width: 0; height: 0; }
        .message-scroll-visible { scrollbar-width: thin; -ms-overflow-style: auto; }
        .message-scroll-visible::-webkit-scrollbar { width: 8px; height: 8px; }
        .message-scroll-visible::-webkit-scrollbar-thumb {
          background: rgba(74, 144, 217, 0.35);
          border-radius: 999px;
        }
        .message-scroll-visible::-webkit-scrollbar-track { background: transparent; }
      `}</style>

      <div style={styles.header}>
        <div style={styles.headerLeft}>
          <div style={styles.headerTitle}>
            {headerTitle}
            {isGroupChat && <span style={styles.headerBadge}>群聊</span>}
          </div>
          <div style={styles.headerSub}>
            {peerTyping ? '对方正在输入...' : (isGroupChat ? '' : '')}
          </div>
        </div>
        <button
          type="button"
          style={styles.searchBtn}
          onClick={() => navigate('/messages/search')}
          title="搜索消息"
        >
          搜索
        </button>
      </div>

      <div
        ref={listRef}
        className={showScrollbar ? 'message-scroll message-scroll-visible' : 'message-scroll message-scroll-hidden'}
        style={styles.list}
        onScroll={() => {
          handleScrollState();
          void handleScrollTop();
        }}
      >
        {loadingHistory && <div style={styles.loadingMore}>加载中...</div>}

        {messages.map((msg: MessageVo, index: number) => {
          const isMe = msg.mine || msg.senderId === currentUserIdVal;
          const otherAvatar = msg.sendAvatar || getSenderAvatar(msg.senderId);
          const senderAvatar = isMe ? currentAvatar : otherAvatar;
          const senderName = isMe ? currentNickname : msg.sendName || nameCache.get(msg.senderId) || (isGroupChat ? '' : '对方');
          const bubbleStyle = isMe ? styles.bubbleMe : styles.bubbleOther;
          const rowStyle = isMe ? styles.msgRowMe : styles.msgRowOther;
          const isFileMsg = msg.msgType === 2 && msg.fileUrl;
          const isImageMsg = isFileMsg && isImageFile(msg.fileName || '');
          const isPdfMsg = isFileMsg && isPdfFile(msg.fileName || '');
          const messageKey = msg.msgId ? `msg-${msg.msgId}` : msg.localId ? `local-${msg.localId}` : `idx-${index}`;

          return (
            <div key={messageKey} style={{ ...styles.msgRow, ...rowStyle }}>
              {/* 头像 - 所有人均可点击跳转 */}
              <button
                type="button"
                style={isMe ? styles.avatarButtonMe : styles.avatarButtonOther}
                onClick={() => navigate(`/profile?userId=${msg.senderId || currentUserIdVal}`)}
                title="查看个人详情"
              >
                <AvatarView src={resolveAvatarUrl(senderAvatar)} alt={senderName} size={34} radius={12} />
              </button>

              <div style={{ ...styles.bubble, ...bubbleStyle }}>
                {/* 群聊显示发送者名称 */}
                {isGroupChat && !isMe && senderName && (
                  <div style={styles.sender}>{senderName}</div>
                )}

                {/* 图片消息 */}
                {isImageMsg ? (
                  <button
                    type="button"
                    style={styles.imageCard}
                    onClick={() => handleDownload(msg)}
                    title="点击查看原图或下载"
                  >
                    <img
                      src={resolveFileUrl(msg.fileUrl)}
                      alt={msg.fileName || '图片'}
                      style={styles.imagePreview}
                      onClick={(e) => {
                        e.stopPropagation();
                        window.open(resolveFileUrl(msg.fileUrl), '_blank');
                      }}
                    />
                  </button>
                ) : isPdfMsg ? (
                  <button
                    type="button"
                    style={styles.fileCard}
                    onClick={() => window.open(resolveFileUrl(msg.fileUrl), '_blank')}
                    title="查看PDF"
                  >
                    <div style={styles.fileIcon}>📄</div>
                    <div style={styles.fileInfo}>
                      <div style={{ ...styles.fileName, ...(isMe ? { color: '#fff' } : {}) }}>
                        {msg.fileName || '文档.pdf'}
                      </div>
                      <div style={{ ...styles.fileSize, ...(isMe ? { color: 'rgba(255,255,255,0.7)' } : {}) }}>
                        {formatFileSize(msg.fileSize)}
                      </div>
                    </div>
                    <div style={{ ...styles.downloadBtn, ...(isMe ? { borderColor: 'rgba(255,255,255,0.4)', color: '#fff' } : {}) }}>
                      预览
                    </div>
                  </button>
                ) : isFileMsg ? (
                  <button
                    type="button"
                    style={styles.fileCard}
                    onClick={() => handleDownload(msg)}
                    title="下载附件"
                  >
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
                  </button>
                ) : (
                  <div style={{ ...styles.content, ...(isMe ? styles.contentMe : {}) }}>
                    {replaceEmojiShortcuts(msg.content || '')}
                  </div>
                )}

                {/* 消息状态与时间 */}
                <div style={styles.metaRow}>
                  {msg.status && isMe && (
                    <span style={styles.status}>
                      {msg.status === 'sending' ? '发送中' : msg.status === 'failed' ? '发送失败' : ''}
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
    padding: '10px 16px',
    borderBottom: '1px solid #dbe7fb',
    background: 'rgba(255,255,255,0.92)',
    color: '#1f2d3d',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  headerLeft: {},
  headerTitle: { fontWeight: 700, fontSize: 15, display: 'flex', alignItems: 'center', gap: 6 },
  headerBadge: {
    fontSize: 10, padding: '1px 6px', borderRadius: 4, background: '#eef3fb', color: '#4a90d9', fontWeight: 600,
  },
  headerSub: { marginTop: 2, fontSize: 12, color: '#7a869a', minHeight: 16 },
  searchBtn: {
    height: 32, padding: '0 12px', borderRadius: 8, border: '1px solid #dbe7fb',
    background: '#f7fbff', color: '#4a90d9', fontWeight: 600, fontSize: 12,
    cursor: 'pointer', flexShrink: 0,
  },
  loadingMore: { textAlign: 'center', padding: 8, fontSize: 12, color: '#999' },
  list: {
    flex: 1,
    overflowY: 'auto',
    padding: '14px 16px',
    display: 'flex',
    flexDirection: 'column',
    gap: 10,
    overscrollBehavior: 'contain',
  },
  empty: { flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#999', fontSize: 16 },
  msgRow: { display: 'flex', gap: 8, alignItems: 'flex-end' },
  msgRowMe: { justifyContent: 'flex-end' },
  msgRowOther: { justifyContent: 'flex-start' },
  avatarButtonMe: { order: 2, border: 'none', background: 'transparent', padding: 0, cursor: 'pointer' },
  avatarButtonOther: { order: 0, border: 'none', background: 'transparent', padding: 0, cursor: 'pointer' },
  bubble: { maxWidth: '66%', padding: '8px 12px', borderRadius: 16, boxShadow: '0 10px 20px rgba(31, 45, 61, 0.06)' },
  bubbleMe: { background: '#1f6feb', color: '#fff', borderBottomRightRadius: 6, order: 1 },
  bubbleOther: { background: '#fff', color: '#1f2d3d', border: '1px solid #dbe7fb', borderBottomLeftRadius: 6 },
  sender: { fontSize: 11, marginBottom: 3, opacity: 0.78, fontWeight: 500 },
  content: { fontSize: 14, color: '#333', lineHeight: 1.5, wordBreak: 'break-word', whiteSpace: 'pre-wrap' },
  contentMe: { color: '#fff' },
  metaRow: { marginTop: 4, display: 'flex', alignItems: 'center', justifyContent: 'flex-end', gap: 8 },
  status: { fontSize: 11, opacity: 0.82 },
  time: { fontSize: 10, opacity: 0.72, textAlign: 'right' },
  fileCard: {
    display: 'flex', alignItems: 'center', gap: 10, cursor: 'pointer',
    padding: '4px 0', border: 'none', background: 'transparent', textAlign: 'left', width: '100%',
  },
  imageCard: {
    display: 'grid', gap: 4, cursor: 'pointer', padding: 0,
    border: 'none', background: 'transparent', textAlign: 'left', width: '100%',
  },
  imagePreview: {
    maxWidth: 260, maxHeight: 260, width: '100%', height: 'auto',
    borderRadius: 12, objectFit: 'cover', background: '#eef5ff',
    border: '1px solid rgba(255,255,255,0.18)',
  },
  fileIcon: { fontSize: 28, flexShrink: 0 },
  fileInfo: { flex: 1, minWidth: 0 },
  fileName: { fontSize: 13, fontWeight: 600, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', color: '#333' },
  fileSize: { fontSize: 11, marginTop: 2, color: '#7a869a' },
  downloadBtn: {
    fontSize: 12, padding: '4px 10px', borderRadius: 8,
    borderWidth: 1, borderStyle: 'solid', borderColor: '#dbe7fb',
    background: 'transparent', cursor: 'pointer', flexShrink: 0, color: '#4a90d9',
  },
};
