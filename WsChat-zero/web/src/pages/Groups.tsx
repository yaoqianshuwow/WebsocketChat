import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '@/api/client';
import AvatarView from '@/components/AvatarView';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import { useMobile } from '@/hooks/useMobile';
import { resolveAvatarUrl } from '@/utils/avatar';
import type { GroupInfoResp, MemberVo } from '@/types';

export default function Groups() {
  const navigate = useNavigate();
  const { createSession } = useChatStore();
  const userInfo = useAuthStore((s) => s.userInfo);
  const isMobile = useMobile();
  const fileRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [groups, setGroups] = useState<GroupInfoResp[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<GroupInfoResp | null>(null);
  const [members, setMembers] = useState<MemberVo[]>([]);
  const [showList, setShowList] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [createName, setCreateName] = useState('');
  const [quickOpen, setQuickOpen] = useState(false);
  const [showEdit, setShowEdit] = useState(false);
  const [editName, setEditName] = useState('');
  const [editNotice, setEditNotice] = useState('');
  const [editAddMode, setEditAddMode] = useState(1);
  const [confirmAction, setConfirmAction] = useState<{
    title: string;
    message: string;
    onConfirm: () => Promise<void>;
  } | null>(null);

  useEffect(() => {
    void loadGroups();
  }, []);

  const loadGroups = async () => {
    const resp = await api.loadMyGroups();
    if (resp.code === 0) {
      const list = resp.data || [];
      setGroups(list);
      if (!isMobile && !selectedGroup && list.length) {
        void selectGroup(list[0]);
      } else if (selectedGroup) {
        const updated = list.find((g) => g.group_id === selectedGroup.group_id);
        if (updated) {
          setSelectedGroup(updated);
          void selectGroup(updated);
        } else {
          setSelectedGroup(null);
          setMembers([]);
        }
      }
    } else {
      alert(resp.message);
    }
  };

  const selectGroup = async (group: GroupInfoResp) => {
    setSelectedGroup(group);
    if (isMobile) setShowList(false);
    const resp = await api.getGroupMemberList(group.group_id || 0);
    if (resp.code === 0) {
      setMembers(resp.memberList || []);
    } else {
      alert(resp.message);
    }
  };

  const backToList = () => setShowList(true);

  const openGroupChat = async (group: GroupInfoResp) => {
    const ok = await createSession(group.group_id || 0, 2, group.name || '');
    if (ok) navigate('/chat');
  };

  const handleCreate = async () => {
    if (!createName.trim()) return;
    const resp = await api.createGroup(createName.trim(), []);
    if (resp.code === 0) {
      setCreateName('');
      setShowCreate(false);
      await loadGroups();
    } else {
      alert(resp.message);
    }
  };

  const goToGroupSearch = () => {
    setQuickOpen(false);
    navigate('/groups/search');
  };

  const handleGroupAvatarUpload = async (file: File) => {
    if (!file || !selectedGroup?.group_id) return;
    setUploading(true);
    try {
      const resp = await api.uploadFile(file);
      if (resp.code === 0 && resp.fileUrl) {
        await api.updateGroupInfo({
          groupId: selectedGroup.group_id,
          avatar: resp.fileUrl,
        });
        await loadGroups();
      } else {
        alert(resp.message);
      }
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = '';
    }
  };

  const openEdit = (group: GroupInfoResp) => {
    setEditName(group.name || '');
    setEditNotice(group.notice || '');
    setEditAddMode(group.add_mode || 1);
    setShowEdit(true);
  };

  const handleEdit = async () => {
    if (!selectedGroup?.group_id) return;
    const resp = await api.updateGroupInfo({
      groupId: selectedGroup.group_id,
      name: editName,
      avatar: selectedGroup.avatar,
      notice: editNotice,
      addMode: editAddMode,
    });
    if (resp.code === 0) {
      setShowEdit(false);
      await loadGroups();
    } else {
      alert(resp.message);
    }
  };

  const handleSetStatus = async () => {
    if (!selectedGroup?.group_id) return;
    const resp = await api.setGroupsStatus(selectedGroup.group_id);
    if (resp.code === 0) {
      alert('操作成功');
      await loadGroups();
    } else {
      alert(resp.message);
    }
  };

  const handleLeave = async () => {
    if (!selectedGroup?.group_id) return;
    const resp = await api.leaveGroup(selectedGroup.group_id);
    if (resp.code === 0) {
      setConfirmAction(null);
      setSelectedGroup(null);
      setMembers([]);
      await loadGroups();
    } else {
      alert(resp.message);
    }
  };

  const handleDismiss = async () => {
    if (!selectedGroup?.group_id) return;
    const resp = await api.dismissGroup(selectedGroup.group_id);
    if (resp.code === 0) {
      setConfirmAction(null);
      setSelectedGroup(null);
      setMembers([]);
      await loadGroups();
    } else {
      alert(resp.message);
    }
  };

  const handleDelete = async () => {
    if (!selectedGroup?.group_id) return;
    const resp = await api.deleteGroups([selectedGroup.group_id]);
    if (resp.code === 0) {
      setConfirmAction(null);
      setSelectedGroup(null);
      setMembers([]);
      await loadGroups();
    } else {
      alert(resp.message);
    }
  };

  const isOwner = userInfo?.user_id === selectedGroup?.owner_id;
  const showConfirm = (title: string, message: string, onConfirm: () => Promise<void>) => {
    setConfirmAction({ title, message, onConfirm });
  };

  const roleLabel = (role: number) => {
    if (role === 2) return '群主';
    if (role === 1) return '管理员';
    return '成员';
  };

  const showSidebar = !isMobile || (isMobile && showList);
  const showDetailPanel = !isMobile || (isMobile && !showList && selectedGroup);

  return (
    <div style={{ ...styles.container, ...(isMobile ? mStyles.container : {}) }}>
      {showSidebar && (
      <div style={{ ...styles.sidebar, ...(isMobile ? mStyles.sidebar : {}) }}>
        <div style={{ ...styles.sideTabs, ...(isMobile ? mStyles.sideTabs : {}) }}>
          <button type="button" style={{ ...styles.sideTab, ...styles.sideTabActive }}>我的群组</button>
          <div style={styles.plusWrap}>
            <button type="button" style={{ ...styles.plusBtn, ...(quickOpen ? styles.plusBtnActive : {}) }} onClick={() => setQuickOpen((v) => !v)} title="更多操作">+</button>
            {quickOpen && (
              <div style={{ ...styles.quickPanel, ...(isMobile ? mStyles.quickPanel : {}) }}>
                <button type="button" style={styles.quickItem} onClick={() => { setQuickOpen(false); setShowCreate(true); }}>创建群组</button>
                <button type="button" style={styles.quickItem} onClick={goToGroupSearch}>搜索群聊</button>
              </div>
            )}
          </div>
        </div>
        <div style={styles.sideHeader}><span style={{ fontSize: 13, color: '#888' }}>{groups.length} 个群组</span></div>
        <div style={styles.groupList}>
          {groups.length === 0 && <div style={styles.empty}>暂无群组</div>}
          {groups.map((group) => {
            const groupName = group.name || `群组 ${group.group_id}`;
            return (
              <button key={group.group_id} type="button" style={{ ...styles.groupItem, ...(selectedGroup?.group_id === group.group_id ? styles.groupItemActive : {}) }} onClick={() => selectGroup(group)} onDoubleClick={() => openGroupChat(group)}>
                <AvatarView src={resolveAvatarUrl(group.avatar)} alt={groupName} size={36} radius={8} style={styles.avatarImg} />
                <div style={styles.itemInfo}>
                  <div style={styles.itemName}>{groupName}</div>
                  <div style={styles.itemMeta}>{group.member_count || 0} 人</div>
                </div>
                {isMobile && <span style={{ fontSize: 12, color: '#4a90d9' }}>&gt;</span>}
              </button>
            );
          })}
        </div>
      </div>
      )}

      {showDetailPanel && (
      <div style={{ ...styles.detail, ...(isMobile ? mStyles.detail : {}) }}>
        {selectedGroup ? (
          <>
            <div style={{ ...styles.detailHeader, ...(isMobile ? mStyles.detailHeader : {}) }}>
              {isMobile && <button type="button" style={mStyles.backBtn} onClick={backToList}>&larr;</button>}
              <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                <button type="button" style={{ border: 'none', background: 'transparent', padding: 0, cursor: 'pointer', flexShrink: 0 }} onClick={() => navigate(`/groups?groupId=${selectedGroup.group_id}`)} title="群详情">
                  <AvatarView src={resolveAvatarUrl(selectedGroup.avatar)} alt={selectedGroup.name || '群'} size={48} radius={12} />
                </button>
                <div>
                  <h3 style={{ ...styles.groupName, ...(isMobile ? mStyles.groupName : {}) }}>{selectedGroup.name || '未命名群组'}</h3>
                  <div style={{ ...styles.detailMeta, ...(isMobile ? mStyles.detailMeta : {}) }}>
                    群 ID：{selectedGroup.group_id} · {selectedGroup.member_count || 0} 人
                    {isOwner ? ' · 你是群主' : ''}
                  </div>
                </div>
              </div>
              <div style={{ ...styles.detailActions, ...(isMobile ? mStyles.detailActions : {}) }}>
                <button type="button" style={{ ...styles.actionBtn, ...(isMobile ? mStyles.actionBtn : {}) }} onClick={() => openGroupChat(selectedGroup)}>进入群聊</button>
                {isOwner && <button type="button" style={{ ...styles.actionBtn, background: '#fa8c16', borderColor: '#fa8c16', ...(isMobile ? mStyles.actionBtn : {}) }} onClick={() => openEdit(selectedGroup)}>编辑</button>}
              </div>
            </div>

            {selectedGroup.notice && <div style={styles.notice}><span style={styles.noticeLabel}>公告：</span>{selectedGroup.notice}</div>}

            <div style={styles.infoSection}>
              <div style={styles.sectionTitle}>群组信息</div>
              <div style={styles.infoGrid}>
                <div style={styles.infoItem}><span style={styles.infoLabel}>加群方式</span><span style={styles.infoValue}>{selectedGroup.add_mode === 0 ? '允许任何人' : selectedGroup.add_mode === 1 ? '需要验证' : '禁止加群'}</span></div>
                <div style={styles.infoItem}><span style={styles.infoLabel}>状态</span><span style={styles.infoValue}>{selectedGroup.status === 0 ? '正常' : '已冻结'}</span></div>
              </div>
            </div>

            <div style={styles.manageSection}>
              <div style={styles.sectionTitle}>管理操作</div>
              {isOwner ? (
                <div style={styles.manageActions}>
                  <button type="button" style={styles.manageBtn} onClick={handleSetStatus}>切换状态</button>
                  <button type="button" style={styles.manageBtn} onClick={() => showConfirm('解散群组', '确定要解散该群组吗？此操作不可撤销。', handleDismiss)}>解散群组</button>
                  <button type="button" style={{ ...styles.manageBtn, background: '#ff4d4f', color: '#fff' }} onClick={() => showConfirm('删除群组', '确定要删除该群组吗？', handleDelete)}>删除群组</button>
                </div>
              ) : (
                <button type="button" style={{ ...styles.manageBtn, background: '#ff4d4f', color: '#fff' }} onClick={() => showConfirm('退出群组', '确定要退出该群组吗？', handleLeave)}>退出群组</button>
              )}
            </div>

            <div style={styles.memberSection}>
              <div style={styles.sectionTitle}>成员列表（{members.length} 人）</div>
              {members.map((m) => (
                <div key={m.userId} style={styles.memberItem}>
                  <button type="button" style={styles.memberAvatarButton} onClick={() => navigate(`/profile?userId=${m.userId}`)} title="查看个人详情">
                    <AvatarView src={resolveAvatarUrl(m.avatar)} alt={m.nickname || `用户 ${m.userId}`} size={32} radius={6} style={styles.memberAvatar} />
                  </button>
                  <div style={styles.memberInfo}>
                    <span style={styles.memberName}>{m.nickname || `用户 ${m.userId}`}</span>
                    <span style={styles.memberId}>ID: {m.userId}</span>
                  </div>
                  <span style={{ ...styles.memberRole, ...(m.role === 2 ? { background: '#fff7e6', color: '#d46b08', border: '1px solid #ffd591' } : {}) }}>{roleLabel(m.role)}</span>
                </div>
              ))}
            </div>
          </>
        ) : (
          <div style={styles.empty}>请选择一个群组查看详情</div>
        )}
      </div>
      )}

      {showCreate && (
        <div style={styles.modal} onClick={() => setShowCreate(false)}>
          <div style={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <h3 style={styles.modalTitle}>创建群组</h3>
            <p style={styles.modalHint}>创建后你就是群主，可在群聊中邀请好友加入</p>
            <input style={styles.modalInput} value={createName} onChange={(e) => setCreateName(e.target.value)} placeholder="输入群组名称" onKeyDown={(e) => e.key === 'Enter' && handleCreate()} />
            <div style={styles.modalActions}>
              <button type="button" style={styles.modalBtnPrimary} onClick={handleCreate}>创建</button>
              <button type="button" style={styles.modalBtnDefault} onClick={() => setShowCreate(false)}>取消</button>
            </div>
          </div>
        </div>
      )}

      {showEdit && (
        <div style={styles.modal} onClick={() => setShowEdit(false)}>
          <div style={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <h3 style={styles.modalTitle}>编辑群组</h3>
            {isOwner && (
              <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 12 }}>
                <img src={resolveAvatarUrl(selectedGroup?.avatar)} alt="群头像" style={{ width: 64, height: 64, borderRadius: 12, objectFit: 'cover', background: '#e6f0ff', cursor: 'pointer', border: '1px solid #dbe7fb' }} onClick={() => fileRef.current?.click()} title="点击修改群头像" />
                <div>
                  <button type="button" style={{ padding: '6px 12px', borderRadius: 8, border: '1px solid #dbe7fb', background: '#f7fbff', cursor: 'pointer', fontSize: 12, color: '#4a90d9' }} onClick={() => fileRef.current?.click()}>{uploading ? '上传中...' : '更换群头像'}</button>
                  <div style={{ fontSize: 11, color: '#999', marginTop: 4 }}>点击头像或按钮上传</div>
                </div>
                <input ref={fileRef} type="file" accept="image/*" style={{ display: 'none' }} onChange={(e) => { const f = e.target.files?.[0]; if (f) handleGroupAvatarUpload(f); }} />
              </div>
            )}
            <label style={styles.modalLabel}>群组名称</label>
            <input style={styles.modalInput} value={editName} onChange={(e) => setEditName(e.target.value)} placeholder="群组名称" />
            <label style={styles.modalLabel}>群公告</label>
            <textarea style={{ ...styles.modalInput, minHeight: 72, resize: 'vertical', fontFamily: 'inherit' }} value={editNotice} onChange={(e) => setEditNotice(e.target.value)} placeholder="群公告" />
            <label style={styles.modalLabel}>加群方式</label>
            <select style={styles.modalInput} value={editAddMode} onChange={(e) => setEditAddMode(Number(e.target.value))}>
              <option value={0}>允许任何人</option>
              <option value={1}>需要验证</option>
              <option value={2}>禁止加群</option>
            </select>
            <div style={styles.modalActions}>
              <button type="button" style={styles.modalBtnPrimary} onClick={handleEdit}>保存</button>
              <button type="button" style={styles.modalBtnDefault} onClick={() => setShowEdit(false)}>取消</button>
            </div>
          </div>
        </div>
      )}

      {confirmAction && (
        <div style={styles.modal} onClick={() => setConfirmAction(null)}>
          <div style={styles.confirmBox} onClick={(e) => e.stopPropagation()}>
            <h4 style={styles.confirmTitle}>{confirmAction.title}</h4>
            <p style={styles.confirmMessage}>{confirmAction.message}</p>
            <div style={styles.confirmActions}>
              <button type="button" style={styles.confirmOk} onClick={async () => { await confirmAction.onConfirm(); }}>确定</button>
              <button type="button" style={styles.confirmCancel} onClick={() => setConfirmAction(null)}>取消</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: { flex: 1, display: 'flex', background: '#fff', overflow: 'hidden' },
  sidebar: { width: 280, borderRight: '1px solid #e8e8e8', display: 'flex', flexDirection: 'column', flexShrink: 0 },
  sideTabs: { display: 'grid', gridTemplateColumns: '1fr 48px', alignItems: 'stretch', borderBottom: '1px solid #e8e8e8' },
  sideTab: { padding: '10px 0', border: 'none', background: '#fafafa', cursor: 'pointer', fontSize: 13, color: '#666', borderBottom: '2px solid transparent' },
  sideTabActive: { background: '#fff', color: '#4a90d9', fontWeight: 600, borderBottomColor: '#4a90d9' },
  plusWrap: { position: 'relative' },
  plusBtn: { width: '100%', height: '100%', border: 'none', background: '#fff', cursor: 'pointer', fontSize: 24, lineHeight: 1, color: '#4a90d9', borderBottom: '2px solid transparent' },
  plusBtnActive: { background: '#eaf3ff', borderBottomColor: '#4a90d9' },
  quickPanel: { position: 'absolute', top: 'calc(100% + 8px)', right: 8, minWidth: 150, padding: 8, borderRadius: 12, background: '#fff', border: '1px solid #dbe7fb', boxShadow: '0 16px 32px rgba(31, 64, 122, 0.14)', zIndex: 20 },
  quickItem: { width: '100%', minHeight: 40, padding: '0 12px', border: 'none', borderRadius: 8, background: '#fff', color: '#1f2d3d', cursor: 'pointer', display: 'flex', alignItems: 'center', fontSize: 13, textAlign: 'left' },
  sideHeader: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '8px 12px' },
  groupList: { flex: 1, overflow: 'auto' },
  groupItem: { display: 'flex', alignItems: 'center', gap: 10, width: '100%', padding: '10px 12px', cursor: 'pointer', border: 'none', borderBottom: '1px solid #f0f0f0', background: '#fff', textAlign: 'left' },
  groupItemActive: { background: '#eaf3ff' },
  avatarImg: { width: 36, height: 36, borderRadius: 8, objectFit: 'cover', flexShrink: 0, background: '#e6f0ff' },
  itemInfo: { flex: 1, minWidth: 0 },
  itemName: { fontSize: 14, fontWeight: 500, color: '#333' },
  itemMeta: { fontSize: 12, color: '#999', marginTop: 2 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
  detail: { flex: 1, padding: 16, overflow: 'auto' },
  detailHeader: { display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 12, marginBottom: 12 },
  groupName: { margin: 0, fontSize: 18, color: '#333' },
  detailMeta: { fontSize: 13, color: '#999', marginTop: 6 },
  detailActions: { display: 'flex', gap: 6, flexShrink: 0 },
  actionBtn: { padding: '6px 14px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 13 },
  notice: { background: '#fff7e6', border: '1px solid #ffd591', padding: '10px 14px', borderRadius: 8, marginBottom: 16, fontSize: 13, color: '#d46b08', lineHeight: 1.6 },
  noticeLabel: { fontWeight: 600 },
  infoSection: { marginBottom: 16 },
  sectionTitle: { fontSize: 14, fontWeight: 600, marginBottom: 8, color: '#333', paddingBottom: 4, borderBottom: '1px solid #f0f0f0' },
  infoGrid: { display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 },
  infoItem: { padding: '8px 12px', background: '#fafafa', borderRadius: 6 },
  infoLabel: { display: 'block', fontSize: 11, color: '#999', marginBottom: 2 },
  infoValue: { fontSize: 13, color: '#333', fontWeight: 500 },
  manageSection: { marginBottom: 16 },
  manageActions: { display: 'flex', gap: 8, flexWrap: 'wrap' as const },
  manageBtn: { padding: '6px 14px', background: '#f5f5f5', color: '#333', border: '1px solid #e0e0e0', borderRadius: 4, cursor: 'pointer', fontSize: 13 },
  memberSection: {},
  memberItem: { display: 'flex', alignItems: 'center', gap: 10, padding: '8px 0', borderBottom: '1px solid #f8f8f8' },
  memberAvatarButton: { border: 'none', background: 'transparent', padding: 0, cursor: 'pointer', flexShrink: 0 },
  memberAvatar: { width: 32, height: 32, borderRadius: 6, objectFit: 'cover', background: '#e6f0ff' },
  memberInfo: { flex: 1, display: 'flex', flexDirection: 'column' },
  memberName: { fontSize: 14, color: '#333', fontWeight: 500 },
  memberId: { fontSize: 11, color: '#999', marginTop: 1 },
  memberRole: { fontSize: 11, padding: '2px 8px', borderRadius: 4, background: '#f5f5f5', color: '#666', flexShrink: 0 },
  modal: { position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 },
  modalContent: { background: '#fff', padding: 24, borderRadius: 12, width: 380, maxWidth: '90vw' },
  modalTitle: { margin: '0 0 4px', fontSize: 16, fontWeight: 600 },
  modalHint: { fontSize: 12, color: '#999', marginBottom: 12 },
  modalLabel: { display: 'block', fontSize: 12, color: '#666', marginBottom: 4, marginTop: 8 },
  modalInput: { width: '100%', padding: '8px 12px', borderRadius: 6, border: '1px solid #e0e0e0', fontSize: 14, marginBottom: 8, boxSizing: 'border-box', outline: 'none' },
  modalActions: { display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 12 },
  modalBtnPrimary: { padding: '8px 20px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  modalBtnDefault: { padding: '8px 20px', background: '#f5f5f5', color: '#666', border: '1px solid #e0e0e0', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  confirmBox: { background: '#fff', borderRadius: 12, width: 320, maxWidth: '90vw', padding: '24px' },
  confirmTitle: { margin: '0 0 12px', fontSize: 16, fontWeight: 600, color: '#333' },
  confirmMessage: { fontSize: 14, color: '#666', marginBottom: 20, lineHeight: 1.5 },
  confirmActions: { display: 'flex', gap: 8, justifyContent: 'flex-end' },
  confirmOk: { padding: '8px 20px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
  confirmCancel: { padding: '8px 20px', background: '#f5f5f5', color: '#666', border: '1px solid #e0e0e0', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
};

const mStyles: Record<string, React.CSSProperties> = {
  container: { flexDirection: 'column' as const },
  sidebar: { width: '100%', borderRight: 'none', flex: 1 },
  sideTabs: { gridTemplateColumns: '1fr 44px' },
  quickPanel: { right: 4, minWidth: 130 },
  detail: { padding: 12, overflow: 'auto', flex: 1 },
  detailHeader: { flexDirection: 'column' as const, alignItems: 'stretch', gap: 8 },
  detailActions: { justifyContent: 'flex-end' as const },
  actionBtn: { fontSize: 12, padding: '5px 10px' },
  groupName: { fontSize: 16 },
  detailMeta: { fontSize: 12 },
  backBtn: { alignSelf: 'flex-start', padding: '4px 10px', border: '1px solid #e0e0e0', borderRadius: 6, background: '#f5f5f5', cursor: 'pointer', fontSize: 14, color: '#333', marginBottom: 4 },
};
