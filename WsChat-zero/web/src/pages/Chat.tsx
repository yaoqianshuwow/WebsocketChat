import SessionList from '@/components/SessionList';
import MessageList from '@/components/MessageList';
import MessageInput from '@/components/MessageInput';
import { useMobile } from '@/hooks/useMobile';
import { useChatStore } from '@/store/chat';

export default function Chat() {
  const isMobile = useMobile();
  const { currentSession, setCurrentSession } = useChatStore();

  return (
    <div style={{ ...styles.container, ...(isMobile ? styles.containerMobile : {}) }}>
      {(!isMobile || !currentSession) && <SessionList />}
      {(!isMobile || currentSession) && (
        <div style={{ ...styles.chatArea, ...(isMobile ? styles.chatAreaMobile : {}) }}>
          {isMobile && currentSession && (
            <div style={styles.mobileChatHeader}>
              <button type="button" style={styles.mobileBackBtn} onClick={() => setCurrentSession(null)}>
                返回会话
              </button>
            </div>
          )}
          <MessageList />
          <MessageInput />
        </div>
      )}
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: 'flex',
    flex: 1,
    height: '100%',
    minHeight: 0,
    overflow: 'hidden',
    padding: 12,
    gap: 12,
    width: '100%',
    boxSizing: 'border-box',
  },
  chatArea: {
    flex: 1,
    height: '100%',
    minWidth: 0,
    minHeight: 0,
    display: 'flex',
    flexDirection: 'column',
    borderRadius: 24,
    overflow: 'hidden',
    border: '1px solid #dbe7fb',
    boxShadow: '0 18px 45px rgba(31, 64, 122, 0.08)',
    background: 'rgba(255,255,255,0.75)',
  },
  containerMobile: {
    flexDirection: 'column',
    padding: 8,
    gap: 0,
  },
  chatAreaMobile: {
    height: '100%',
    flex: 1,
    borderRadius: 18,
  },
  mobileChatHeader: {
    padding: '10px 12px',
    borderBottom: '1px solid #dbe7fb',
    background: 'rgba(255,255,255,0.96)',
  },
  mobileBackBtn: {
    height: 34,
    padding: '0 12px',
    borderRadius: 10,
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
    color: '#4a90d9',
    fontSize: 13,
    fontWeight: 700,
  },
};
