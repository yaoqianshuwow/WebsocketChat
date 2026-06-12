import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '@/api/client';
import { useChatStore } from '@/store/chat';
import type { GroupInfoResp, MemberVo } from '@/types';

export default function Groups() {
  const navigate = useNavigate();
  const { createSession } = useChatStore();
  const [groups, setGroups] = useState<GroupInfoResp[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<GroupInfoResp | null>(null);
  const [members, setMembers] = useState<MemberVo[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [groupName, setGroupName] = useState('');
  const [memberIds, setMemberIds] = useState('');

  useEffect(() => {
    void loadGroups();
  }, []);

  const loadGroups = async () => {
    const resp = await api.loadMyGroups();
    if (resp.code === 0) {
      setGroups(resp.data || []);
      if (!selectedGroup && resp.data?.length) {
        void selectGroup(resp.data[0]);
      }
    } else {
      alert(resp.message);
    }
  };

  const selectGroup = async (group: GroupInfoResp) => {
    setSelectedGroup(group);
    const resp = await api.getGroupMemberList(group.group_id || 0);
    if (resp.code === 0) {
      setMembers(resp.memberList || []);
    } else {
      alert(resp.message);
    }
  };

  const openGroupChat = async (group: GroupInfoResp) => {
    const ok = await createSession(group.group_id || 0, 2, group.name || '');
    if (ok) {
      navigate('/chat');
    }
  };

  const handleCreate = async () => {
    if (!groupName.trim()) return;
    const ids = memberIds
      .split(',')
      .map((item) => Number(item.trim()))
      .filter((item) => Number.isFinite(item) && item > 0);

    const resp = await api.createGroup(groupName.trim(), ids);
    if (resp.code === 0) {
      setGroupName('');
      setMemberIds('');
      setShowCreate(false);
      await loadGroups();
    } else {
      alert(resp.message);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.sidebar}>
        <div style={styles.header}>
          <span>我的群组</span>
          <button type="button" style={styles.createBtn} onClick={() => setShowCreate(true)}>
            + 创建
          </button>
        </div>

        {groups.map((group) => (
          <button
            key={group.group_id}
            type="button"
            style={{ ...styles.item, ...(selectedGroup?.group_id === group.group_id ? styles.itemActive : {}) }}
            onClick={() => selectGroup(group)}
            onDoubleClick={() => openGroupChat(group)}
          >
            <div style={styles.avatar}>{(group.name || `G${group.group_id}`).charAt(0)}</div>
            <div style={styles.itemInfo}>
              <div style={styles.itemName}>{group.name || `群组 ${group.group_id}`}</div>
              <div style={styles.itemMeta}>{group.member_count || 0} 人</div>
            </div>
            <span style={styles.enter}>进入</span>
          </button>
        ))}
        {groups.length === 0 && <div style={styles.empty}>暂无群组</div>}
      </div>

      <div style={styles.detail}>
        {selectedGroup ? (
          <>
            <div style={styles.detailHeader}>
              <div>
                <h3 style={styles.groupName}>{selectedGroup.name}</h3>
                <div style={styles.detailMeta}>群 ID：{selectedGroup.group_id} · {selectedGroup.member_count || 0} 人</div>
              </div>
              <button type="button" style={styles.openBtn} onClick={() => openGroupChat(selectedGroup)}>
                进入群聊
              </button>
            </div>
            {selectedGroup.notice && <div style={styles.notice}>公告：{selectedGroup.notice}</div>}
            <div style={styles.memberTitle}>成员列表</div>
            {members.map((m) => (
              <div key={m.userId} style={styles.memberItem}>
                <div style={styles.avatar}>{(m.nickname || `U${m.userId}`).charAt(0)}</div>
                <span style={styles.memberName}>{m.nickname || `用户 ${m.userId}`}</span>
                <span style={styles.memberRole}>{m.role === 2 ? '群主' : m.role === 1 ? '管理员' : '成员'}</span>
              </div>
            ))}
          </>
        ) : (
          <div style={styles.empty}>请选择一个群组查看详情</div>
        )}
      </div>

      {showCreate && (
        <div style={styles.modal}>
          <div style={styles.modalContent}>
            <h3 style={styles.modalTitle}>创建群组</h3>
            <input style={styles.modalInput} value={groupName} onChange={(e) => setGroupName(e.target.value)} placeholder="群组名称" />
            <input style={styles.modalInput} value={memberIds} onChange={(e) => setMemberIds(e.target.value)} placeholder="成员 ID，逗号分隔" />
            <div style={styles.modalActions}>
              <button type="button" style={styles.modalBtn} onClick={handleCreate}>
                创建
              </button>
              <button type="button" style={{ ...styles.modalBtn, background: '#999' }} onClick={() => setShowCreate(false)}>
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
  container: { flex: 1, display: 'flex', background: '#fff', overflow: 'hidden' },
  sidebar: { width: 280, borderRight: '1px solid #e8e8e8', display: 'flex', flexDirection: 'column' },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '12px 16px',
    borderBottom: '1px solid #e8e8e8',
    fontWeight: 600,
  },
  createBtn: { padding: '4px 12px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  item: {
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    width: '100%',
    padding: '10px 16px',
    cursor: 'pointer',
    border: 'none',
    borderBottom: '1px solid #f0f0f0',
    background: '#fff',
    textAlign: 'left',
  },
  itemActive: { background: '#eaf3ff' },
  avatar: {
    width: 36,
    height: 36,
    borderRadius: 8,
    background: '#4a90d9',
    color: '#fff',
    display: 'grid',
    placeItems: 'center',
    fontWeight: 600,
    flexShrink: 0,
  },
  itemInfo: { flex: 1, minWidth: 0 },
  itemName: { fontSize: 14, fontWeight: 500, color: '#333' },
  itemMeta: { fontSize: 12, color: '#999', marginTop: 2 },
  enter: { fontSize: 12, color: '#4a90d9', flexShrink: 0 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
  detail: { flex: 1, padding: 16, overflow: 'auto' },
  detailHeader: { display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 12, marginBottom: 16 },
  groupName: { margin: 0, fontSize: 18, color: '#333' },
  detailMeta: { fontSize: 13, color: '#999', marginTop: 6 },
  openBtn: { padding: '6px 14px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 13 },
  notice: { background: '#fff7e6', border: '1px solid #ffd591', padding: '8px 12px', borderRadius: 6, marginBottom: 16, fontSize: 13, color: '#d46b08' },
  memberTitle: { fontSize: 14, fontWeight: 600, marginBottom: 8, color: '#333' },
  memberItem: { display: 'flex', alignItems: 'center', gap: 10, padding: '8px 0' },
  memberName: { fontSize: 14, color: '#333', flex: 1 },
  memberRole: { fontSize: 12, color: '#999' },
  modal: {
    position: 'fixed',
    inset: 0,
    background: 'rgba(0,0,0,0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modalContent: { background: '#fff', padding: 24, borderRadius: 12, width: 360 },
  modalTitle: { marginTop: 0 },
  modalInput: {
    width: '100%',
    padding: '8px 12px',
    borderRadius: 6,
    border: '1px solid #e0e0e0',
    fontSize: 14,
    marginBottom: 12,
    boxSizing: 'border-box',
    outline: 'none',
  },
  modalActions: { display: 'flex', gap: 8, justifyContent: 'flex-end' },
  modalBtn: { padding: '8px 20px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
};
