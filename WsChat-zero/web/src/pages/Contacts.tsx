import { useEffect, useMemo, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '@/api/client';
import AvatarView from '@/components/AvatarView';
import { useChatStore } from '@/store/chat';
import { useMobile } from '@/hooks/useMobile';
import { resolveAvatarUrl } from '@/utils/avatar';
import type { ApplyVo, ContactVo } from '@/types';

type TabKey = 'list';

function safeText(value: string | undefined, fallback: string) {
  if (!value) return fallback;
  if (value.includes('�')) return fallback;
  return value;
}

function displayName(name?: string, fallback = '用户') {
  return safeText(name, fallback);
}

export default function Contacts() {
  const { contacts, loadContacts, createSession } = useChatStore();
  const navigate = useNavigate();
  const isMobile = useMobile();
  const [applyList, setApplyList] = useState<ApplyVo[]>([]);
  const [tab, setTab] = useState<TabKey>('list');
  const [applyOpen, setApplyOpen] = useState(false);
  const [quickOpen, setQuickOpen] = useState(false);

  // 详情弹窗
  const [detailContact, setDetailContact] = useState<ContactVo | null>(null);
  const [detailInfo, setDetailInfo] = useState<ContactVo[] | null>(null);

  // 确认对话框
  const [confirmAction, setConfirmAction] = useState<{
    title: string;
    message: string;
    onConfirm: () => Promise<void>;
  } | null>(null);

  useEffect(() => {
    void refresh();
  }, []);

  const refresh = useCallback(async () => {
    await Promise.all([loadContacts(), loadApplies()]);
  }, [loadContacts]);

  const loadApplies = useCallback(async () => {
    const resp = await api.getApplyList();
    if (resp.code === 0) {
      setApplyList(resp.data || []);
    }
  }, []);

  const openChat = async (contactId: number, sessionName?: string) => {
    const ok = await createSession(contactId, 1, sessionName || '');
    if (ok) {
      navigate('/chat');
    }
  };

  const handleApply = async (toId: number) => {
    const resp = await api.applyContact(toId, '你好，想加个好友');
    alert(resp.message);
    if (resp.code === 0) {
      await loadApplies();
    }
  };

  const handlePass = async (applyId: number, status: number) => {
    const resp = await api.passContactApply(applyId, status);
    alert(resp.message);
    if (resp.code === 0) {
      await refresh();
    }
  };

  const handleDeleteContact = async (contactId: number) => {
    const resp = await api.deleteContact(contactId);
    alert(resp.message);
    if (resp.code === 0) {
      await refresh();
      setDetailContact(null);
    }
  };

  const handleBlackContact = async (contactId: number) => {
    const resp = await api.blackContact(contactId);
    alert(resp.message);
    if (resp.code === 0) {
      await refresh();
      setDetailContact(null);
    }
  };

  const showConfirm = (title: string, message: string, onConfirm: () => Promise<void>) => {
    setConfirmAction({ title, message, onConfirm });
  };

  const openDetail = async (contact: ContactVo) => {
    setDetailContact(contact);
    const resp = await api.getContactInfo2(contact.contactId);
    if (resp.code === 0) {
      // getContactInfo 返回的是数组
      setDetailInfo(resp.data || []);
    } else {
      setDetailInfo(null);
    }
  };

  const applyCount = useMemo(() => applyList.filter((item) => item.status === 0).length, [applyList]);

  return (
    <div style={{ ...styles.container, ...(isMobile ? mStyles.container : {}) }}>
      {/* 顶部栏 */}
      <div style={{ ...styles.tabs, ...(isMobile ? mStyles.tabs : {}) }}>
        <button
          type="button"
          style={{ ...styles.tab, ...(tab === 'list' ? styles.tabActive : {}), ...(isMobile ? mStyles.tab : {}) }}
          onClick={() => setTab('list')}
        >
          联系人
        </button>
        <div style={styles.plusWrap}>
          <button
            type="button"
            style={{ ...styles.plusBtn, ...(applyOpen || quickOpen ? styles.plusBtnActive : {}), ...(isMobile ? mStyles.plusBtn : {}) }}
            onClick={() => {
              setQuickOpen((v) => !v);
              setApplyOpen(false);
            }}
            title="更多操作"
          >
            +
            {applyCount > 0 && <span style={styles.badge}>{applyCount}</span>}
          </button>
          {quickOpen && (
            <div style={{ ...styles.quickPanel, ...(isMobile ? mStyles.quickPanel : {}) }}>
              <button
                type="button"
                style={styles.quickItem}
                onClick={() => {
                  setApplyOpen(true);
                  setTab('list');
                  setQuickOpen(false);
                }}
              >
                <span>新朋友</span>
                {applyCount > 0 && <span style={styles.quickBadge}>{applyCount}</span>}
              </button>
              <button
                type="button"
                style={styles.quickItem}
                onClick={() => {
                  navigate('/contacts/search');
                  setApplyOpen(false);
                  setQuickOpen(false);
                }}
              >
                <span>查找用户</span>
              </button>
            </div>
          )}
        </div>
      </div>

      {/* 好友申请面板 */}
      {applyOpen && (
        <div style={styles.applyPanel}>
          <div style={styles.applyHeader}>
            <span>新朋友</span>
            <button type="button" style={styles.linkBtn} onClick={refresh}>
              刷新
            </button>
          </div>
          <div style={styles.applyList}>
            {applyList.length === 0 && <div style={styles.emptyMini}>暂无新的好友申请</div>}
            {applyList.map((item) => {
              const name = displayName(item.nickname, `用户 ${item.fromId}`);
              return (
                <div key={item.applyId} style={styles.applyItem}>
                  <button
                    type="button"
                    style={styles.avatarButton}
                    onClick={() => navigate(`/profile?userId=${item.fromId}`)}
                    title="查看个人详情"
                  >
                    <AvatarView src={resolveAvatarUrl(item.avatar)} alt={name} size={36} radius={8} style={styles.avatarImg} />
                  </button>
                  <div style={styles.itemInfo}>
                    <div style={styles.itemName}>{name}</div>
                    <div style={styles.itemStatus}>{safeText(item.remark, '好友申请')}</div>
                  </div>
                  {item.status === 0 ? (
                    <div style={styles.actions}>
                      <button type="button" style={styles.passBtn} onClick={() => handlePass(item.applyId, 1)}>
                        同意
                      </button>
                      <button type="button" style={styles.rejectBtn} onClick={() => handlePass(item.applyId, 2)}>
                        拒绝
                      </button>
                    </div>
                  ) : (
                    <span style={styles.doneText}>{item.status === 1 ? '已同意' : '已拒绝'}</span>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* 主体内容 */}
      <div style={styles.content}>
        {tab === 'list' && (
          <div>
            {contacts.length === 0 && <div style={styles.empty}>暂无联系人</div>}
            {contacts.map((contact) => {
              const name = displayName(contact.nickname, '好友');
              return (
                <div key={contact.contactId} style={{ ...styles.item, ...(isMobile ? mStyles.item : {}) }}>
                  <div
                    style={{ ...styles.itemBody, ...(isMobile ? mStyles.itemBody : {}) }}
                    onClick={() => openChat(contact.contactId, contact.nickname)}
                  >
                    <AvatarView
                      src={resolveAvatarUrl(contact.avatar)}
                      alt={name}
                      size={isMobile ? 32 : 36}
                      radius={8}
                      style={{ ...styles.avatarImg, ...(isMobile ? mStyles.avatarImg : {}) }}
                    />
                    <div style={styles.itemInfo}>
                      <div style={{ ...styles.itemName, ...(isMobile ? mStyles.itemName : {}) }}>{name}</div>
                      <div style={styles.itemStatus}>
                        {contact.status === 0 ? '在线' : '离线'}
                      </div>
                    </div>
                  </div>
                  <div style={{ ...styles.itemActions, ...(isMobile ? mStyles.itemActions : {}) }}>
                    <button
                      type="button"
                      style={{ ...styles.actionSmall, ...(isMobile ? mStyles.actionSmall : {}) }}
                      onClick={(e) => { e.stopPropagation(); openDetail(contact); }}
                      title="详情"
                    >
                      {isMobile ? '···' : '详情'}
                    </button>
                    <button
                      type="button"
                      style={{ ...styles.actionSmall, ...(isMobile ? mStyles.actionSmall : {}) }}
                      onClick={(e) => {
                        e.stopPropagation();
                        openChat(contact.contactId, contact.nickname);
                      }}
                      title="聊天"
                    >
                      聊天
                    </button>
                    <button
                      type="button"
                      style={{ ...styles.actionSmall, color: '#ff4d4f', ...(isMobile ? mStyles.actionSmall : {}) }}
                      onClick={(e) => {
                        e.stopPropagation();
                        showConfirm('删除联系人', `确定要删除联系人「${name}」吗？`, () => handleDeleteContact(contact.contactId));
                      }}
                      title="删除"
                    >
                      删除
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* 联系人详情弹窗 */}
      {detailContact && (
        <div style={styles.modal} onClick={() => setDetailContact(null)}>
          <div style={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <div style={styles.modalHeader}>
              <h3 style={styles.modalTitle}>联系人详情</h3>
              <button
                type="button"
                style={styles.modalClose}
                onClick={() => setDetailContact(null)}
              >
                ✕
              </button>
            </div>
            <div style={styles.modalBody}>
              <div style={styles.detailRow}>
                <span style={styles.detailLabel}>联系人 ID</span>
                <span style={styles.detailValue}>{detailContact.contactId}</span>
              </div>
              <div style={styles.detailRow}>
                <span style={styles.detailLabel}>昵称</span>
                <span style={styles.detailValue}>{displayName(detailContact.nickname, '未设置')}</span>
              </div>
              <div style={styles.detailRow}>
                <span style={styles.detailLabel}>状态</span>
                <span style={styles.detailValue}>{detailContact.status === 0 ? '在线' : '离线'}</span>
              </div>
              <div style={styles.detailRow}>
                <span style={styles.detailLabel}>类型</span>
                <span style={styles.detailValue}>{detailContact.contactType === 1 ? '好友' : '其他'}</span>
              </div>
              {detailInfo && detailInfo.length > 0 && (
                <>
                  <div style={styles.detailDivider} />
                  <div style={styles.detailRow}>
                    <span style={styles.detailLabel}>更多信息</span>
                    <span style={styles.detailValue}>查询成功</span>
                  </div>
                </>
              )}
            </div>
            <div style={styles.modalFooter}>
              <button
                type="button"
                style={styles.modalBtnPrimary}
                onClick={() => {
                  openChat(detailContact.contactId, detailContact.nickname);
                  setDetailContact(null);
                }}
              >
                发送消息
              </button>
              <button
                type="button"
                style={{ ...styles.modalBtnDanger }}
                onClick={() => {
                  showConfirm('拉黑联系人', `确定要拉黑「${displayName(detailContact.nickname, '该用户')}」吗？拉黑后将不再接收对方消息。`, async () => {
                    await handleBlackContact(detailContact.contactId);
                  });
                }}
              >
                拉黑
              </button>
              <button
                type="button"
                style={styles.modalBtnDefault}
                onClick={() => setDetailContact(null)}
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 确认对话框 */}
      {confirmAction && (
        <div style={styles.modal} onClick={() => setConfirmAction(null)}>
          <div style={{ ...styles.confirmBox }} onClick={(e) => e.stopPropagation()}>
            <h4 style={styles.confirmTitle}>{confirmAction.title}</h4>
            <p style={styles.confirmMessage}>{confirmAction.message}</p>
            <div style={styles.confirmActions}>
              <button
                type="button"
                style={styles.confirmOk}
                onClick={async () => {
                  await confirmAction.onConfirm();
                  setConfirmAction(null);
                }}
              >
                确定
              </button>
              <button
                type="button"
                style={styles.confirmCancel}
                onClick={() => setConfirmAction(null)}
              >
                取消
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { flex: 1, display: 'flex', flexDirection: 'column', background: '#fff' },
  tabs: { display: 'grid', gridTemplateColumns: '1fr 64px', alignItems: 'stretch', borderBottom: '1px solid #e8e8e8' },
  tab: {
    padding: '12px 0',
    border: 'none',
    background: '#fafafa',
    cursor: 'pointer',
    fontSize: 14,
    color: '#666',
    borderBottom: '2px solid transparent',
  },
  tabActive: { background: '#fff', color: '#4a90d9', fontWeight: 600, borderBottomColor: '#4a90d9' },
  plusWrap: { position: 'relative' },
  plusBtn: {
    position: 'relative',
    width: '100%',
    height: '100%',
    border: 'none',
    background: '#fff',
    cursor: 'pointer',
    fontSize: 28,
    lineHeight: 1,
    color: '#4a90d9',
    borderBottom: '2px solid transparent',
  },
  plusBtnActive: { background: '#eaf3ff', borderBottomColor: '#4a90d9' },
  quickPanel: {
    position: 'absolute',
    top: 'calc(100% + 8px)',
    right: 8,
    minWidth: 150,
    padding: 8,
    borderRadius: 12,
    background: '#fff',
    border: '1px solid #dbe7fb',
    boxShadow: '0 16px 32px rgba(31, 64, 122, 0.14)',
    zIndex: 20,
  },
  quickItem: {
    width: '100%',
    minHeight: 40,
    padding: '0 12px',
    border: 'none',
    borderRadius: 8,
    background: '#fff',
    color: '#1f2d3d',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    fontSize: 13,
    textAlign: 'left',
  },
  quickBadge: {
    minWidth: 18,
    height: 18,
    padding: '0 5px',
    borderRadius: 9,
    background: '#f5222d',
    color: '#fff',
    fontSize: 11,
    display: 'grid',
    placeItems: 'center',
  },
  badge: {
    position: 'absolute',
    top: 8,
    right: 7,
    minWidth: 18,
    height: 18,
    padding: '0 5px',
    borderRadius: 9,
    background: '#f5222d',
    color: '#fff',
    fontSize: 11,
    display: 'grid',
    placeItems: 'center',
  },
  applyPanel: { borderBottom: '1px solid #e8e8e8', background: '#f7fbff' },
  applyHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '10px 12px',
    fontSize: 13,
    color: '#333',
    borderBottom: '1px solid #edf2f7',
  },
  linkBtn: { border: 'none', background: 'transparent', color: '#4a90d9', cursor: 'pointer', padding: 0 },
  applyList: { maxHeight: 220, overflow: 'auto' },
  applyItem: {
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    padding: '10px 12px',
    borderBottom: '1px solid #f0f0f0',
    background: '#fff',
  },
  avatarButton: { border: 'none', background: 'transparent', padding: 0, cursor: 'pointer', flexShrink: 0 },
  content: { flex: 1, overflow: 'auto', padding: 8 },
  item: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    padding: '6px 12px',
    borderBottom: '1px solid #f0f0f0',
    background: '#fff',
  },
  itemBody: {
    flex: 1,
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    cursor: 'pointer',
    minWidth: 0,
    border: 'none',
    background: 'transparent',
    padding: '4px 0',
    textAlign: 'left',
  },
  avatarImg: { width: 36, height: 36, borderRadius: 8, objectFit: 'cover', flexShrink: 0, background: '#e6f0ff' },
  itemInfo: { flex: 1, minWidth: 0 },
  itemName: { fontSize: 14, fontWeight: 500, color: '#333' },
  itemStatus: { fontSize: 12, color: '#999', marginTop: 2, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' },
  itemActions: { display: 'flex', gap: 4, flexShrink: 0 },
  actionSmall: {
    padding: '4px 8px',
    border: '1px solid #e0e0e0',
    borderRadius: 4,
    background: '#fff',
    color: '#4a90d9',
    cursor: 'pointer',
    fontSize: 11,
    whiteSpace: 'nowrap',
  },
  actions: { display: 'flex', gap: 6 },
  passBtn: { padding: '4px 12px', background: '#52c41a', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  rejectBtn: { padding: '4px 12px', background: '#ff4d4f', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  doneText: { fontSize: 12, color: '#999' },
  addBtn: { padding: '4px 12px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
  emptyMini: { padding: 20, textAlign: 'center', color: '#999', fontSize: 13 },
  searchBar: { display: 'flex', gap: 8, padding: '8px 12px' },
  searchInput: { flex: 1, padding: '8px 12px', borderRadius: 6, border: '1px solid #e0e0e0', fontSize: 14, outline: 'none' },
  searchBtn: { padding: '8px 16px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  // 弹窗
  modal: {
    position: 'fixed',
    inset: 0,
    background: 'rgba(0,0,0,0.45)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modalContent: { background: '#fff', borderRadius: 12, width: 380, maxWidth: '90vw', overflow: 'hidden' },
  modalHeader: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '16px 20px', borderBottom: '1px solid #f0f0f0' },
  modalTitle: { margin: 0, fontSize: 16, fontWeight: 600 },
  modalClose: { border: 'none', background: 'transparent', fontSize: 18, color: '#999', cursor: 'pointer', padding: 4 },
  modalBody: { padding: '16px 20px' },
  detailRow: { display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid #f8f8f8' },
  detailLabel: { fontSize: 13, color: '#888' },
  detailValue: { fontSize: 13, color: '#333', fontWeight: 500 },
  detailDivider: { height: 1, background: '#eee', margin: '8px 0' },
  modalFooter: { display: 'flex', gap: 8, padding: '12px 20px', borderTop: '1px solid #f0f0f0', justifyContent: 'flex-end' },
  modalBtnPrimary: { padding: '6px 16px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 13 },
  modalBtnDanger: { padding: '6px 16px', background: '#ff4d4f', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 13 },
  modalBtnDefault: { padding: '6px 16px', background: '#f5f5f5', color: '#666', border: '1px solid #e0e0e0', borderRadius: 6, cursor: 'pointer', fontSize: 13 },
  // 确认对话框
  confirmBox: { background: '#fff', borderRadius: 12, width: 320, maxWidth: '90vw', padding: '24px' },
  confirmTitle: { margin: '0 0 12px', fontSize: 16, fontWeight: 600, color: '#333' },
  confirmMessage: { fontSize: 14, color: '#666', marginBottom: 20, lineHeight: 1.5 },
  confirmActions: { display: 'flex', gap: 8, justifyContent: 'flex-end' },
  confirmOk: { padding: '8px 20px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  confirmCancel: { padding: '8px 20px', background: '#f5f5f5', color: '#666', border: '1px solid #e0e0e0', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
};

// ── 移动端适配 ──
const mStyles: Record<string, React.CSSProperties> = {
  container: { padding: 0 },
  tabs: { gridTemplateColumns: '1fr 48px' },
  tab: { fontSize: 13, padding: '10px 0' },
  plusBtn: { fontSize: 24 },
  quickPanel: { right: 4, minWidth: 130 },
  item: { padding: '4px 8px', gap: 6, flexWrap: 'nowrap' },
  itemBody: { gap: 8, padding: '4px 0' },
  avatarImg: { width: 32, height: 32 },
  itemName: { fontSize: 13 },
  itemActions: { gap: 2, flexDirection: 'column' as const },
  actionSmall: { padding: '2px 6px', fontSize: 10 },
  addBtn: { padding: '3px 8px', fontSize: 11 },
};
