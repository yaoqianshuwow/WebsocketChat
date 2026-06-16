import { useMemo, useState } from 'react';
import { useMobile } from '@/hooks/useMobile';
import { getPptAssistantUrl, PPT_ASSISTANT_NAME } from '@/utils/pptAssistant';

export default function PptAssistant() {
  const isMobile = useMobile();
  const [reloadKey, setReloadKey] = useState(0);

  const assistantUrl = useMemo(() => getPptAssistantUrl(), []);

  const openExternal = () => {
    window.open(assistantUrl, '_blank', 'noopener,noreferrer');
  };

  return (
    <div style={{ ...styles.shell, ...(isMobile ? styles.shellMobile : {}) }}>
      <div style={{ ...styles.panel, ...(isMobile ? styles.panelMobile : {}) }}>
        <div style={styles.hero}>
          <div style={styles.avatarWrap}>
            <div style={styles.avatar}>
              <span style={styles.avatarText}>AI</span>
            </div>
          </div>
          <div style={styles.heroText}>
            <div style={styles.title}>{PPT_ASSISTANT_NAME}</div>
            <div style={styles.subTitle}>接入 codingagent 工作台，支持对话生成、预览和下载 PPT</div>
            <div style={styles.meta}>当前地址：{assistantUrl}</div>
          </div>
          <div style={styles.heroActions}>
            <button type="button" style={styles.primaryBtn} onClick={() => setReloadKey((v) => v + 1)}>
              刷新工作台
            </button>
            <button type="button" style={styles.secondaryBtn} onClick={openExternal}>
              新窗口打开
            </button>
          </div>
        </div>

        <div style={styles.tipRow}>
          <span style={styles.tip}>1. 生成 PPT 仍由独立 AI 服务处理</span>
          <span style={styles.tip}>2. WsChat 这里只负责入口和会话承载</span>
          <span style={styles.tip}>3. 如果本地 8123 未启动，可改 `VITE_PPT_ASSISTANT_URL`</span>
        </div>
      </div>

      <div style={{ ...styles.workspace, ...(isMobile ? styles.workspaceMobile : {}) }}>
        <iframe
          key={reloadKey}
          title={PPT_ASSISTANT_NAME}
          src={assistantUrl}
          style={styles.iframe}
        />
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  shell: {
    flex: 1,
    minHeight: 0,
    display: 'grid',
    gridTemplateRows: 'auto 1fr',
    gap: 12,
    padding: 12,
    background: 'linear-gradient(180deg, rgba(245,249,255,1) 0%, rgba(238,244,255,1) 100%)',
  },
  shellMobile: {
    padding: 0,
    gap: 0,
    gridTemplateRows: 'auto 1fr',
  },
  panel: {
    borderRadius: 20,
    padding: 16,
    background: 'rgba(255,255,255,0.94)',
    border: '1px solid #dbe7fb',
    boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)',
  },
  panelMobile: {
    borderRadius: 0,
    borderLeft: 'none',
    borderRight: 'none',
    borderTop: 'none',
    boxShadow: 'none',
  },
  hero: {
    display: 'flex',
    gap: 14,
    alignItems: 'center',
    flexWrap: 'wrap',
  },
  avatarWrap: {
    width: 56,
    height: 56,
    borderRadius: 18,
    overflow: 'hidden',
    flexShrink: 0,
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
  },
  avatar: {
    width: '100%',
    height: '100%',
    display: 'grid',
    placeItems: 'center',
    fontSize: 18,
    fontWeight: 800,
    color: '#4a90d9',
    background: 'linear-gradient(135deg, #eef5ff 0%, #dbe7fb 100%)',
  },
  avatarText: {
    fontSize: 18,
    fontWeight: 800,
    color: '#4a90d9',
  },
  heroText: {
    flex: 1,
    minWidth: 220,
  },
  title: {
    fontSize: 18,
    fontWeight: 800,
    color: '#1f2d3d',
  },
  subTitle: {
    marginTop: 4,
    fontSize: 13,
    color: '#5f6b7a',
    lineHeight: 1.5,
  },
  meta: {
    marginTop: 6,
    fontSize: 12,
    color: '#7a869a',
    wordBreak: 'break-all',
  },
  heroActions: {
    display: 'flex',
    gap: 8,
    flexWrap: 'wrap',
    marginLeft: 'auto',
  },
  primaryBtn: {
    height: 36,
    padding: '0 14px',
    borderRadius: 10,
    border: 'none',
    background: '#4a90d9',
    color: '#fff',
    fontWeight: 700,
    cursor: 'pointer',
  },
  secondaryBtn: {
    height: 36,
    padding: '0 14px',
    borderRadius: 10,
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
    color: '#4a90d9',
    fontWeight: 700,
    cursor: 'pointer',
  },
  tipRow: {
    marginTop: 12,
    display: 'flex',
    gap: 8,
    flexWrap: 'wrap',
  },
  tip: {
    padding: '6px 10px',
    borderRadius: 999,
    background: '#eef5ff',
    color: '#4a90d9',
    fontSize: 12,
    border: '1px solid #dbe7fb',
  },
  workspace: {
    borderRadius: 20,
    overflow: 'hidden',
    border: '1px solid #dbe7fb',
    background: '#fff',
    minHeight: 0,
    boxShadow: '0 16px 40px rgba(31, 64, 122, 0.08)',
  },
  workspaceMobile: {
    borderRadius: 0,
    borderLeft: 'none',
    borderRight: 'none',
    borderBottom: 'none',
    boxShadow: 'none',
  },
  iframe: {
    width: '100%',
    height: '100%',
    minHeight: 'calc(100dvh - 220px)',
    border: 'none',
    display: 'block',
    background: '#fff',
  },
};
