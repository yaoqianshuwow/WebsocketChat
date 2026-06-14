import { useEffect, useMemo, useRef, useState } from 'react';
import api from '@/api/client';
import { useAuthStore } from '@/store/auth';
import { pickAvatar } from '@/utils/avatar';

export default function Profile() {
  const { userInfo, loadUserInfo } = useAuthStore();
  const fileRef = useRef<HTMLInputElement>(null);
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [form, setForm] = useState({
    nickname: '',
    avatar: '',
    sex: '',
    age: '',
    bio: '',
  });

  useEffect(() => {
    void loadUserInfo();
  }, [loadUserInfo]);

  useEffect(() => {
    if (!userInfo) return;
    setForm({
      nickname: userInfo.nickname || '',
      avatar: userInfo.avatar || '',
      sex: userInfo.sex || '',
      age: userInfo.age ? String(userInfo.age) : '',
      bio: userInfo.bio || '',
    });
  }, [userInfo]);

  const avatarSrc = useMemo(() => {
    return pickAvatar(userInfo?.user_id || 0, form.avatar || userInfo?.avatar);
  }, [form.avatar, userInfo?.avatar, userInfo?.user_id]);

  const handleSave = async () => {
    setSaving(true);
    try {
      const resp = await api.updateUserInfo({
        nickname: form.nickname.trim(),
        avatar: form.avatar.trim(),
        sex: form.sex.trim(),
        age: form.age ? Number(form.age) : undefined,
        bio: form.bio.trim(),
      });
      alert(resp.message);
      if (resp.code === 0) {
        await loadUserInfo();
      }
    } finally {
      setSaving(false);
    }
  };

  const handleAvatarUpload = async (file: File) => {
    if (!file) return;
    setUploading(true);
    try {
      const resp = await api.uploadAvatar(file);
      if (resp.code === 0 && resp.fileUrl) {
        setForm((prev) => ({ ...prev, avatar: resp.fileUrl || '' }));
        await loadUserInfo();
        alert('头像上传成功');
      } else {
        alert(resp.message);
      }
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = '';
    }
  };

  return (
    <div style={styles.wrap}>
      <div style={styles.card}>
        <div style={styles.header}>
          <div>
            <div style={styles.title}>个人信息</div>
            <div style={styles.sub}>头像、昵称、性别和简介都在这里维护</div>
          </div>
          <div style={styles.avatarBox}>
            <img src={avatarSrc} alt="avatar" style={styles.avatar} />
            <button type="button" style={styles.avatarBtn} onClick={() => fileRef.current?.click()} disabled={uploading}>
              {uploading ? '上传中' : '更换头像'}
            </button>
            <input
              ref={fileRef}
              type="file"
              accept="image/*"
              style={{ display: 'none' }}
              onChange={(e) => {
                const file = e.target.files?.[0];
                if (file) void handleAvatarUpload(file);
              }}
            />
          </div>
        </div>

        <div style={styles.grid}>
          <label style={styles.field}>
            <span style={styles.label}>用户名</span>
            <input style={styles.input} value={userInfo?.username || ''} readOnly />
          </label>
          <label style={styles.field}>
            <span style={styles.label}>昵称</span>
            <input
              style={styles.input}
              value={form.nickname}
              onChange={(e) => setForm((prev) => ({ ...prev, nickname: e.target.value }))}
              placeholder="请输入昵称"
            />
          </label>
          <label style={styles.field}>
            <span style={styles.label}>性别</span>
            <input
              style={styles.input}
              value={form.sex}
              onChange={(e) => setForm((prev) => ({ ...prev, sex: e.target.value }))}
              placeholder="例如：男 / 女 / 保密"
            />
          </label>
          <label style={styles.field}>
            <span style={styles.label}>年龄</span>
            <input
              style={styles.input}
              value={form.age}
              onChange={(e) => setForm((prev) => ({ ...prev, age: e.target.value }))}
              placeholder="请输入年龄"
              inputMode="numeric"
            />
          </label>
        </div>

        <label style={styles.field}>
          <span style={styles.label}>签名</span>
          <textarea
            style={{ ...styles.input, minHeight: 110, resize: 'vertical' }}
            value={form.bio}
            onChange={(e) => setForm((prev) => ({ ...prev, bio: e.target.value }))}
            placeholder="写点个人介绍吧"
          />
        </label>

        <div style={styles.meta}>
          <span>状态：{userInfo?.status === 0 ? '在线' : '离线'}</span>
          <span>角色：{userInfo?.role === 1 ? '管理员' : '普通用户'}</span>
          <span>ID：{userInfo?.user_id || 0}</span>
        </div>

        <div style={styles.actions}>
          <button type="button" style={styles.primaryBtn} onClick={handleSave} disabled={saving}>
            {saving ? '保存中' : '保存修改'}
          </button>
        </div>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  wrap: {
    flex: 1,
    minHeight: 0,
    padding: 16,
    overflow: 'auto',
    background: 'linear-gradient(180deg, rgba(245,249,255,1) 0%, rgba(238,244,255,1) 100%)',
  },
  card: {
    maxWidth: 960,
    margin: '0 auto',
    borderRadius: 24,
    padding: 24,
    background: 'rgba(255,255,255,0.92)',
    border: '1px solid #dbe7fb',
    boxShadow: '0 18px 45px rgba(31, 64, 122, 0.08)',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    gap: 20,
    alignItems: 'flex-start',
    marginBottom: 20,
  },
  title: { fontSize: 22, fontWeight: 800, color: '#1f2d3d' },
  sub: { marginTop: 6, fontSize: 13, color: '#7a869a' },
  avatarBox: {
    display: 'grid',
    justifyItems: 'center',
    gap: 10,
  },
  avatar: {
    width: 96,
    height: 96,
    borderRadius: 24,
    objectFit: 'cover',
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
  },
  avatarBtn: {
    height: 36,
    padding: '0 14px',
    borderRadius: 10,
    border: '1px solid #dbe7fb',
    background: '#f7fbff',
    color: '#4a90d9',
    fontWeight: 700,
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
    gap: 16,
  },
  field: {
    display: 'grid',
    gap: 8,
    marginTop: 16,
  },
  label: {
    fontSize: 13,
    color: '#5f6b7a',
    fontWeight: 700,
  },
  input: {
    width: '100%',
    borderRadius: 14,
    border: '1px solid #dbe7fb',
    background: '#f8fbff',
    padding: '12px 14px',
    fontSize: 14,
    color: '#1f2d3d',
  },
  meta: {
    marginTop: 16,
    display: 'flex',
    flexWrap: 'wrap',
    gap: 12,
    fontSize: 13,
    color: '#5f6b7a',
  },
  actions: {
    marginTop: 20,
    display: 'flex',
    justifyContent: 'flex-end',
  },
  primaryBtn: {
    height: 42,
    padding: '0 18px',
    borderRadius: 12,
    border: 'none',
    background: '#4a90d9',
    color: '#fff',
    fontWeight: 800,
  },
};
