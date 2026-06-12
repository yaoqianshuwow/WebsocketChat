import SessionList from '@/components/SessionList';
import MessageList from '@/components/MessageList';
import MessageInput from '@/components/MessageInput';

export default function Chat() {
  return (
    <div style={styles.container}>
      <SessionList />
      <div style={styles.chatArea}>
        <MessageList />
        <MessageInput />
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { display: 'flex', flex: 1, overflow: 'hidden' },
  chatArea: { flex: 1, display: 'flex', flexDirection: 'column' },
};
