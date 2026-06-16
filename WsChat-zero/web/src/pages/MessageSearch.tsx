import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '@/api/client';
import { useMobile } from '@/hooks/useMobile';
import { resolveFileUrl } from '@/utils/fileUrl';
import type { MessageVo } from '@/types';

function formatTime(ts: number) {
  const d = new Date(ts * 1000);
  return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
}

export default function MessageSearch() {
  const navigate = useNavigate();
  const isMobile = useMobile();
  const [keyword, setKeyword] = useState('');
  const [senderId, setSenderId] = useState('');
  const [results, setResults] = useState<MessageVo[]>([]);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);

  const handleSearch = async () => {
    const kw = keyword.trim();
    if (!kw && !senderId.trim()) return;
    setLoading(true);
    setSearched(true);
    try {
      const resp = await api.searchMessages({
        keyword: kw || undefined,
        senderId: senderId ? Number(senderId) : undefined,
        page: 1,
        size: 50,
      });
      if (resp.code === 0) {
        setResults(resp.data || []);
      } else {
        alert(resp.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ ...styles.wrap, ...(isMobile ? styles.wrapMobile : {}) }}>
      <div style={styles.header}>
        <button type="button" style={styles.backBtn} onClick={() => navigate(-1)}>
          返回
        </button>
        <div>
          <div style={styles.title}>消息搜索</div>
          <div style={styles.sub}>搜索关键词和发送者，不阻挡消息列表</div>
        </div>
      </div>

      <div style={styles.searchBar}>
        <input
          style={styles.input}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="搜索关键词"
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <input
          style={{ ...styles.input, width: isMobile ? 80 : 120 }}
          value={senderId}
          onChange={(e) => setSenderId(e.target.value)}
          placeholder="发送者ID"
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <button type="button" style={styles.searchBtn} onClick={handleSearch} disabled={loading}>
          {loading ? '搜索中' : '搜索'}
        </button>
      </div>

      <div style={styles.list}>
        {!searched && <div style={styles.empty}>输入关键词后搜索</div>}
        {searched && results.length === 0 && <div style={styles.empty}>未找到匹配消息</div>}
        {results.map((msg, i) => (
          <div key={msg.msgId || i} style={styles.item}>
            <div style={styles.itemHeader}>
              <span style={styles.sender}>发送者: {msg.senderId}</span>
              <span style={styles.time}>{formatTime(msg.createdAt)}</span>
            </div>
            <div style={styles.content}>{msg.content || '[文件消息]'}</div>
            <div style={styles.meta}>
              会话: {msg.receiverId} · 类型: {msg.msgType === 1 ? '文本' : '文件'}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  wrap: { flex: 1, minHeight: 0, padding: 12, display: 'flex', flexDirection: 'column', gap: 12, background: 'linear-gradient(180deg, rgba(245,249,255,1) 0%, rgba(238,244,255,1) 100%)' },
  wrapMobile: { padding: 0, gap: 0 },
  header: { display: 'flex', alignItems: 'center', gap: 12, padding: 16, borderRadius: 18, background: 'rgba(255,255,255,0.94)', border: '1px solid #dbe7fb', boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)' },
  backBtn: { height: 36, padding: '0 12px', borderRadius: 10, border: '1px solid #dbe7fb', background: '#f7fbff', color: '#4a90d9', fontWeight: 700, cursor: 'pointer' },
  title: { fontSize: 18, fontWeight: 800, color: '#1f2d3d' },
  sub: { marginTop: 4, fontSize: 12, color: '#7a869a' },
  searchBar: { display: 'flex', gap: 8, padding: 12, borderRadius: 18, background: 'rgba(255,255,255,0.94)', border: '1px solid #dbe7fb', boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)' },
  input: { flex: 1, minWidth: 0, borderRadius: 12, border: '1px solid #dbe7fb', padding: '10px 12px', fontSize: 14, outline: 'none' },
  searchBtn: { height: 42, padding: '0 16px', borderRadius: 12, border: 'none', background: '#4a90d9', color: '#fff', fontWeight: 700, cursor: 'pointer' },
  list: { flex: 1, minHeight: 0, overflow: 'auto', borderRadius: 18, background: 'rgba(255,255,255,0.94)', border: '1px solid #dbe7fb', boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)' },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
  item: { padding: '12px 14px', borderBottom: '1px solid #f3f5f7' },
  itemHeader: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 6 },
  sender: { fontSize: 12, color: '#4a90d9', fontWeight: 600 },
  time: { fontSize: 11, color: '#a0a7b4' },
  content: { fontSize: 14, color: '#333', lineHeight: 1.5, marginBottom: 4, wordBreak: 'break-word' },
  meta: { fontSize: 11, color: '#bbb' },
};
