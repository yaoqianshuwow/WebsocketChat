import { useEffect, useState } from 'react';
import api from '@/api/client';
import type { GroupInfoResp, MemberVo } from '@/types';

export default function Groups() {
  const [groups, setGroups] = useState<GroupInfoResp[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [groupName, setGroupName] = useState('');
  const [memberIds, setMemberIds] = useState('');
  const [selectedGroup, setSelectedGroup] = useState<GroupInfoResp | null>(null);
  const [members, setMembers] = useState<MemberVo[]>([]);

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = async () => {
    const resp = await api.loadMyGroups();
    if (resp.code === 0) {
      setGroups(resp.data || []);
    }
  };

  const handleCreate = async () => {
    if (!groupName.trim()) return;
    const ids = memberIds.split(',').map(Number).filter(Boolean);
    const resp = await api.createGroup(groupName, ids);
    if (resp.code === 0) {
      setShowCreate(false);
      setGroupName('');
      setMemberIds('');
      loadGroups();
    } else {
      alert(resp.message);
    }
  };

  const handleSelectGroup = async (g: GroupInfoResp) => {
    setSelectedGroup(g);
    const resp = await api.getGroupMemberList(g.group_id!);
    if (resp.code === 0) {
      setMembers(resp.memberList || []);
    }
  };

  return (
    <div style={{ ...styles.container, flexDirection: selectedGroup ? 'row' : 'column' }}>
      <div style={styles.sidebar}>
        <div style={styles.header}>
          <span>我的群组</span>
          <button style={styles.createBtn} onClick={() => setShowCreate(true)}>+ 创建</button>
        </div>

        {groups.map((g) => (
          <div
            key={g.group_id}
            style={{ ...styles.item, ...(selectedGroup?.group_id === g.group_id ? styles.itemActive : {}) }}
            onClick={() => handleSelectGroup(g)}
          >
            <div style={styles.avatar}>{g.name?.charAt(0) || 'G'}</div>
            <div style={styles.itemInfo}>
              <div style={styles.itemName}>{g.name || `群组${g.group_id}`}</div>
              <div style={styles.itemMeta}>{g.member_count || 0} 人</div>
            </div>
          </div>
        ))}
        {groups.length === 0 && <div style={styles.empty}>暂无群组</div>}
      </div>

      {selectedGroup && (
        <div style={styles.detail}>
          <div style={styles.detailHeader}>
            <h3>{selectedGroup.name}</h3>
            <span style={styles.detailMeta}>{selectedGroup.member_count || 0} 人</span>
          </div>
          {selectedGroup.notice && <div style={styles.notice}>公告：{selectedGroup.notice}</div>}
          <div style={styles.memberTitle}>成员列表</div>
          {members.map((m) => (
            <div key={m.userId} style={styles.memberItem}>
              <div style={styles.avatar}>{m.nickname?.charAt(0) || '?'}</div>
              <span style={styles.memberName}>{m.nickname || `用户${m.userId}`}</span>
              <span style={styles.memberRole}>{m.role === 2 ? '群主' : m.role === 1 ? '管理员' : '成员'}</span>
            </div>
          ))}
        </div>
      )}

      {showCreate && (
        <div style={styles.modal}>
          <div style={styles.modalContent}>
            <h3>创建群组</h3>
            <input style={styles.modalInput} value={groupName} onChange={(e) => setGroupName(e.target.value)} placeholder="群组名称" />
            <input style={styles.modalInput} value={memberIds} onChange={(e) => setMemberIds(e.target.value)} placeholder="成员ID（逗号分隔）" />
            <div style={styles.modalActions}>
              <button style={styles.modalBtn} onClick={handleCreate}>创建</button>
              <button style={{ ...styles.modalBtn, background: '#999' }} onClick={() => setShowCreate(false)}>取消</button>
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
  header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '12px 16px', borderBottom: '1px solid #e8e8e8' },
  createBtn: { padding: '4px 12px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 12 },
  item: { display: 'flex', padding: '10px 16px', cursor: 'pointer', borderBottom: '1px solid #f0f0f0' },
  itemActive: { background: '#e6f0ff' },
  avatar: { width: 36, height: 36, borderRadius: 8, background: '#4a90d9', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 600, flexShrink: 0 },
  itemInfo: { marginLeft: 10, flex: 1 },
  itemName: { fontSize: 14, fontWeight: 500, color: '#333' },
  itemMeta: { fontSize: 12, color: '#999', marginTop: 2 },
  empty: { padding: 40, textAlign: 'center', color: '#999', fontSize: 14 },
  detail: { flex: 1, padding: 16, overflow: 'auto' },
  detailHeader: { display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 },
  detailMeta: { fontSize: 13, color: '#999' },
  notice: { background: '#fff7e6', border: '1px solid #ffd591', padding: '8px 12px', borderRadius: 6, marginBottom: 16, fontSize: 13, color: '#d46b08' },
  memberTitle: { fontSize: 14, fontWeight: 600, marginBottom: 8, color: '#333' },
  memberItem: { display: 'flex', alignItems: 'center', gap: 10, padding: '8px 0' },
  memberName: { fontSize: 14, color: '#333', flex: 1 },
  memberRole: { fontSize: 12, color: '#999' },
  modal: { position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 },
  modalContent: { background: '#fff', padding: 24, borderRadius: 12, width: 360 },
  modalInput: { width: '100%', padding: '8px 12px', borderRadius: 6, border: '1px solid #e0e0e0', fontSize: 14, marginBottom: 12, boxSizing: 'border-box', outline: 'none' },
  modalActions: { display: 'flex', gap: 8, justifyContent: 'flex-end' },
  modalBtn: { padding: '8px 20px', background: '#4a90d9', color: '#fff', border: 'none', borderRadius: 6, cursor: 'pointer', fontSize: 14 },
};
