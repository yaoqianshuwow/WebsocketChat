const API_BASE = '/api/v1';

class ApiClient {
  private token: string | null = null;

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

    return res.json();
  }

  private async upload<T>(path: string, file: File, extra?: Record<string, string>) {
    const form = new FormData();
    form.append('file', file);
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
    return res.json() as Promise<T>;
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

  getApplyList() {
    return this.request<import('../types').ApplyListResp>('/contact/getNewContactList', {});
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
  getMessageList(sessionId: number, page = 1, size = 20) {
    return this.request<import('../types').MessageListResp>('/message/getMessageList', { sessionId, page, size });
  }

  getGroupMessageList(groupId: number, page = 1, size = 20) {
    return this.request<import('../types').MessageListResp>('/message/getGroupMessageList', { groupId, page, size });
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

  uploadFile(file: File) {
    return this.upload<import('../types').CommonResp>('/message/uploadFile', file);
  }

  uploadAvatar(file: File) {
    return this.upload<import('../types').CommonResp>('/message/uploadAvatar', file);
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
