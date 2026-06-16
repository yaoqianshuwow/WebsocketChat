const API_BASE = import.meta.env.VITE_API_BASE_URL
  ? `${import.meta.env.VITE_API_BASE_URL}/api/v1`
  : '/api/v1';

class ApiClient {
  private token: string | null = null;

  private handleAuthExpired() {
    this.setToken(null);
    window.dispatchEvent(new Event('wschat-auth-expired'));
  }

  setToken(token: string | null) {
    this.token = token;
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  }

  getToken(): string | null {
    if (!this.token) {
      this.token = localStorage.getItem('token');
    }
    return this.token;
  }

  private async request<T>(path: string, body?: unknown): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    const token = this.getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const res = await fetch(`${API_BASE}${path}`, {
      method: 'POST',
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}: ${res.statusText}`);
    }

    const data = await res.json();
    if (data && typeof data === 'object' && 'code' in data && (data as { code?: number }).code === 401) {
      this.handleAuthExpired();
    }
    return data as T;
  }

  private async upload<T>(path: string, file: File, fieldName = 'file', extra?: Record<string, string>) {
    const form = new FormData();
    form.append(fieldName, file);
    Object.entries(extra || {}).forEach(([key, value]) => form.append(key, value));

    const headers: Record<string, string> = {};
    const token = this.getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const res = await fetch(`${API_BASE}${path}`, {
      method: 'POST',
      headers,
      body: form,
    });

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}: ${res.statusText}`);
    }
    const data = await res.json();
    if (data && typeof data === 'object' && 'code' in data && (data as { code?: number }).code === 401) {
      this.handleAuthExpired();
    }
    return data as T;
  }

  // ── 认证 ──
  login(data: { username: string; password: string }) {
    return this.request<import('../types').LoginResp>('/login', data);
  }

  register(data: { username: string; password: string; phone?: string; nickname?: string }) {
    return this.request<import('../types').RegisterResp>('/register', data);
  }

  // ── 用户 ──
  getUserInfo(userId?: number) {
    return this.request<import('../types').UserInfoResp>('/user/getUserInfo', { userId });
  }

  updateUserInfo(data: import('../types').UpdateUserInfoReq) {
    return this.request<import('../types').CommonResp>('/user/updateUserInfo', data);
  }

  searchUsers(keyword: string, page = 1, size = 20) {
    return this.request<import('../types').SearchUsersResp>('/user/searchUsers', { keyword, page, size });
  }

  // ── 联系人 ──
  getContactList() {
    return this.request<import('../types').ContactListResp>('/contact/getUserList', {});
  }

  getContactInfo(contactId: number) {
    return this.request<import('../types').UserInfoResp>('/contact/getContactInfo', { contactId });
  }

  applyContact(toId: number, remark = '') {
    return this.request<import('../types').CommonResp>('/contact/applyContact', { toId, remark });
  }

  passContactApply(applyId: number, status: number) {
    return this.request<import('../types').CommonResp>('/contact/passContactApply', { applyId, status });
  }

  deleteContact(contactId: number) {
    return this.request<import('../types').CommonResp>('/contact/deleteContact', { contactId });
  }

  blackContact(contactId: number) {
    return this.request<import('../types').CommonResp>('/contact/blackContact', { contactId });
  }

  getApplyList() {
    return this.request<import('../types').ApplyListResp>('/contact/getNewContactList', {});
  }

  getContactInfo2(contactId: number) {
    return this.request<import('../types').ContactListResp>('/contact/getContactInfo', { contactId });
  }

  // ── 会话 ──
  getSessionList(sessionType = 0) {
    return this.request<import('../types').SessionListResp>('/session/getUserSessionList', { sessionType });
  }

  createSession(peerId: number, sessionType: number, sessionName = '') {
    return this.request<import('../types').SessionResp>('/session/openSession', { peerId, sessionType, sessionName });
  }

  deleteSession(sessionId: number) {
    return this.request<import('../types').CommonResp>('/session/deleteSession', { sessionId });
  }

  // ── 消息 ──
  getMessageList(sessionId: number, beforeId = 0, size = 20) {
    return this.request<import('../types').MessageListResp>('/message/getMessageList', { sessionId, beforeId, size });
  }

  getGroupMessageList(groupId: number, beforeId = 0, size = 20) {
    return this.request<import('../types').MessageListResp>('/message/getGroupMessageList', { groupId, beforeId, size });
  }

  searchMessages(data: import('../types').SearchMessagesReq) {
    return this.request<import('../types').MessageListResp>('/message/searchMessages', data);
  }

  getRecentMessages(sessionId: number, limit = 20) {
    return this.request<import('../types').MessageListResp>('/message/getRecentMessages', { sessionId, limit });
  }

  // ── 群组 ──
  createGroup(groupName: string, memberIds: number[]) {
    return this.request<import('../types').GroupInfoResp>('/group/createGroup', { groupName, memberIds });
  }

  getGroupInfo(groupId: number) {
    return this.request<import('../types').GroupInfoResp>('/group/getGroupInfo', { groupId });
  }

  loadMyGroups() {
    return this.request<import('../types').GroupListResp>('/group/loadMyGroup', {});
  }

  getGroupMemberList(groupId: number) {
    return this.request<import('../types').GroupMemberListResp>('/group/getGroupMemberList', { groupId });
  }

  updateGroupInfo(data: import('../types').UpdateGroupInfoReq) {
    return this.request<import('../types').CommonResp>('/group/updateGroupInfo', data);
  }

  searchGroupList(keyword: string) {
    return this.request<import('../types').SearchGroupListResp>('/group/searchGroupList', { keyword });
  }

  setGroupsStatus(groupId: number) {
    return this.request<import('../types').CommonResp>('/group/setGroupsStatus', { groupId });
  }

  dismissGroup(groupId: number) {
    return this.request<import('../types').CommonResp>('/group/dismissGroup', { groupId });
  }

  deleteGroups(ids: number[]) {
    return this.request<import('../types').CommonResp>('/group/deleteGroups', { ids });
  }

  uploadFile(file: File) {
    return this.upload<import('../types').UploadFileResp>('/message/uploadFile', file);
  }

  uploadAvatar(file: File) {
    return this.upload<import('../types').UploadAvatarResp>('/message/uploadAvatar', file, 'avatar');
  }

  downloadFile(fileUrl: string) {
    const token = this.getToken();
    const baseUrl = import.meta.env.VITE_API_BASE_URL || '';
    const a = document.createElement('a');
    a.href = `${baseUrl}/api/v1/message/downloadFile?fileUrl=${encodeURIComponent(fileUrl)}`;
    if (token) a.href += `&token=${encodeURIComponent(token)}`;
    a.target = '_blank';
    a.rel = 'noopener noreferrer';
    a.click();
  }

  joinGroup(groupId: number) {
    return this.request<import('../types').CommonResp>('/group/joinGroup', { groupId });
  }

  leaveGroup(groupId: number) {
    return this.request<import('../types').CommonResp>('/group/leaveGroup', { groupId });
  }

  removeGroupMembers(groupId: number, memberIds: number[]) {
    return this.request<import('../types').CommonResp>('/group/removeGroupMembers', { groupId, memberIds });
  }
}

export const api = new ApiClient();
export default api;
