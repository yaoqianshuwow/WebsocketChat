import { useEffect, useRef, useState } from 'react';
import api from '@/api/client';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import { useMobile } from '@/hooks/useMobile';
import wsClient from '@/ws/client';
import type { MessageVo } from '@/types';

const EMOJIS = ['😀', '😂', '😭', '🥹', '😎', '❤️'];
const TYPING_IDLE_MS = 3000;

export default function MessageInput() {
  const [text, setText] = useState('');
  const [busy, setBusy] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);
  const inputRef = useRef<HTMLTextAreaElement>(null);
  const typingTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const typingSentRef = useRef(false);
  const { currentSession, addMessage, markMessageStatus, wsState } = useChatStore();
  const { user, userInfo } = useAuthStore();
  const isMobile = useMobile();

  const userId = user?.userId ?? user?.user_id ?? userInfo?.user_id ?? 0;
  const session = currentSession;

  useEffect(() => {
    if (!isMobile || !session) return;
    const handleFocus = () => {
      setTimeout(() => {
        inputRef.current?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
      }, 300);
    };
    const el = inputRef.current;
    el?.addEventListener('focus', handleFocus);
    return () => el?.removeEventListener('focus', handleFocus);
  }, [session, isMobile]);

  useEffect(() => {
    return () => {
      if (typingTimerRef.current) clearTimeout(typingTimerRef.current);
    };
  }, []);

  const stopTyping = () => {
    const activeSession = session!;
    if (typingTimerRef.current) {
      clearTimeout(typingTimerRef.current);
      typingTimerRef.current = null;
    }
    if (typingSentRef.current) {
      wsClient.sendTyping(activeSession.peerId, activeSession.sessionType, activeSession.sessionId, false);
      typingSentRef.current = false;
    }
  };

  useEffect(() => {
    return () => {
      stopTyping();
    };
  }, [session?.sessionId]);

  if (!session) {
    return null;
  }

  const scheduleTyping = () => {
    wsClient.sendTyping(session.peerId, session.sessionType, session.sessionId, true);
    typingSentRef.current = true;

    if (typingTimerRef.current) {
      clearTimeout(typingTimerRef.current);
    }
    typingTimerRef.current = setTimeout(() => {
      stopTyping();
    }, TYPING_IDLE_MS);
  };

  const sendText = () => {
    const content = text.trim();
    if (!content) return;

    stopTyping();

    const localId = `local-${Date.now()}`;
    const localMessage: MessageVo = {
      localId,
      senderId: userId,
      receiverId: session.peerId,
      msgType: 1,
      content,
      createdAt: Math.floor(Date.now() / 1000),
      status: 'sending',
      mine: true,
    };

    addMessage(localMessage);
    const sent = wsClient.sendMessage(content, session.peerId, session.sessionType, session.sessionId);
    markMessageStatus(localId, sent ? 'sent' : 'failed');
    setText('');
  };

  const handleSend = () => {
    if (busy || wsState === 'connecting' || wsState === 'reconnecting') return;
    sendText();
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleUpload = async (file: File) => {
    if (!file) return;
    setBusy(true);
    try {
      const resp = await api.uploadFile(file);
      if (resp.code !== 0 || !resp.fileUrl) {
        alert(resp.message);
        return;
      }

      stopTyping();

      const localId = `local-${Date.now()}`;
      const localMessage: MessageVo = {
        localId,
        senderId: userId,
        receiverId: session.peerId,
        msgType: 2,
        content: '',
        fileUrl: resp.fileUrl,
        fileName: resp.fileName || file.name,
        fileSize: resp.fileSize || file.size,
        createdAt: Math.floor(Date.now() / 1000),
        status: 'sending',
        mine: true,
      };
      addMessage(localMessage);

      const sent = wsClient.sendMessage(
        '',
        session.peerId,
        session.sessionType,
        session.sessionId,
        2,
        {
          fileUrl: resp.fileUrl,
          fileName: resp.fileName || file.name,
          fileSize: resp.fileSize || file.size,
        },
      );
      markMessageStatus(localId, sent ? 'sent' : 'failed');
    } finally {
      setBusy(false);
      if (fileRef.current) fileRef.current.value = '';
    }
  };

  const insertEmoji = (emoji: string) => {
    setText((prev) => `${prev}${emoji}`);
    setTimeout(scheduleTyping, 0);
  };

  return (
    <div style={styles.container}>
      <div style={styles.toolbar}>
        <div style={styles.toolbarLeft}>
          <button type="button" style={styles.toolBtn} onClick={() => fileRef.current?.click()}>
            附件
          </button>
          <input
            ref={fileRef}
            type="file"
            style={{ display: 'none' }}
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) void handleUpload(file);
            }}
          />
          <div style={styles.emojiBar}>
            {EMOJIS.map((emoji) => (
              <button key={emoji} type="button" style={styles.emojiBtn} onClick={() => insertEmoji(emoji)} title={`发送 ${emoji}`}>
                {emoji}
              </button>
            ))}
          </div>
        </div>
        <div style={styles.hint}>
          {busy ? '正在上传' : wsState === 'connected' ? '连接正常，可以直接发送' : wsState === 'reconnecting' ? '连接波动，正在重连' : '等待连接建立'}
        </div>
      </div>

      <div style={styles.editor}>
        <textarea
          ref={inputRef}
          style={{ ...styles.input, ...(isMobile ? styles.inputMobile : {}) }}
          value={text}
          onChange={(e) => {
            const value = e.target.value;
            setText(value);
            if (value.trim()) {
              if (typingTimerRef.current) clearTimeout(typingTimerRef.current);
              scheduleTyping();
            } else {
              stopTyping();
            }
          }}
          onKeyDown={handleKeyDown}
          onBlur={stopTyping}
          placeholder="输入消息，Enter 发送，Shift + Enter 换行"
          rows={isMobile ? 2 : 3}
        />
        <button type="button" style={styles.sendBtn} onClick={handleSend} disabled={!text.trim() || busy}>
          发送
        </button>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    padding: 14,
    borderTop: '1px solid #dbe7fb',
    background: 'rgba(255,255,255,0.96)',
    display: 'grid',
    gap: 10,
  },
  toolbar: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 12 },
  toolbarLeft: { display: 'flex', alignItems: 'center', gap: 10, flexWrap: 'wrap' },
  toolBtn: {
    height: 34,
    padding: '0 12px',
    borderRadius: 10,
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
    color: '#4a90d9',
    cursor: 'pointer',
  },
  emojiBar: { display: 'flex', gap: 6, flexWrap: 'wrap' },
  emojiBtn: {
    minWidth: 34,
    height: 34,
    padding: '0 8px',
    borderRadius: 10,
    border: '1px solid #e5ebf7',
    background: '#fff',
    cursor: 'pointer',
    fontSize: 16,
  },
  hint: { fontSize: 12, color: '#999', textAlign: 'right' },
  editor: { display: 'flex', alignItems: 'flex-end', gap: 10 },
  input: {
    flex: 1,
    resize: 'none',
    minHeight: 84,
    maxHeight: 140,
    padding: '12px 14px',
    borderRadius: 16,
    border: '1px solid #dbe7fb',
    background: '#f8fbff',
    color: '#333',
    lineHeight: 1.6,
    outline: 'none',
  },
  inputMobile: {
    minHeight: 56,
    maxHeight: 100,
    fontSize: 16,
  },
  sendBtn: {
    minWidth: 88,
    height: 44,
    padding: '0 18px',
    borderRadius: 14,
    border: 'none',
    color: '#fff',
    fontWeight: 700,
    background: '#4a90d9',
    cursor: 'pointer',
  },
};
