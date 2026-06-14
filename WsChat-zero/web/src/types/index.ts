// ── 通用 ──
export interface CommonResp {
  code: number;
  message: string;
}

export interface UploadFileResp {
  code: number;
  message: string;
  fileUrl?: string;
  fileName?: string;
  fileSize?: number;
}

// ── 用户 ──
export interface LoginReq {
  username: string;
  password: string;
}

export interface LoginResp {
  code: number;
  message: string;
  token: string;
  user_id: number;
  nickname: string;
  avatar: string;
}

export interface RegisterReq {
  username: string;
  password: string;
  phone?: string;
  nickname?: string;
}

export interface RegisterResp {
  code: number;
  message: string;
  user_id: number;
}

export interface SessionResp {
  code: number;
  message: string;
  sessionId?: number;
  peerId?: number;
  sessionType?: number;
  sessionName?: string;
}

export interface UserInfoResp {
  code: number;
  message: string;
  user_id?: number;
  username?: string;
  nickname?: string;
  avatar?: string;
  sex?: string;
  age?: number;
  bio?: string;
  phone?: string;
  status?: number;
  role?: number;
}

export interface SearchUsersResp {
  code: number;
  message: string;
  data: UserInfoResp[];
  total: number;
}

export interface UpdateUserInfoReq {
  nickname?: string;
  avatar?: string;
  sex?: string;
  age?: number;
  bio?: string;
}

export interface UploadAvatarResp {
  code: number;
  message: string;
  fileUrl?: string;
  fileName?: string;
  fileSize?: number;
}

// ── 消息 ──
export interface MessageVo {
  msgId?: number;
  localId?: string;
  senderId: number;
  receiverId: number;
  msgType: number;
  content?: string;
  fileUrl?: string;
  fileName?: string;
  fileSize?: number;
  createdAt: number;
  status?: 'sending' | 'sent' | 'failed';
  mine?: boolean;
  sendName?: string;
  sendAvatar?: string;
}

export interface MessageListResp {
  code: number;
  message: string;
  data: MessageVo[];
  total: number;
}

// ── 联系人 ──
export interface ContactVo {
  contactId: number;
  contactType: number;
  nickname?: string;
  avatar?: string;
  status: number;
}

export interface ContactListResp {
  code: number;
  message: string;
  data: ContactVo[];
}

// ── 申请 ──
export interface ApplyVo {
  applyId: number;
  fromId: number;
  nickname?: string;
  remark?: string;
  status: number;
}

export interface ApplyListResp {
  code: number;
  message: string;
  data: ApplyVo[];
}

// ── 会话 ──
export interface SessionVo {
  sessionId: number;
  peerId: number;
  sessionType: number;
  sessionName: string;
  lastMsgContent?: string;
  lastMsgTime?: number;
  unreadCount: number;
  avatar?: string;
}

export interface SessionListResp {
  code: number;
  message: string;
  data: SessionVo[];
}

// ── 群组 ──
export interface GroupInfoResp {
  code: number;
  message: string;
  group_id?: number;
  name?: string;
  avatar?: string;
  owner_id?: number;
  member_count?: number;
  add_mode?: number;
  status?: number;
  notice?: string;
}

export interface UpdateGroupInfoReq {
  groupId: number;
  name?: string;
  avatar?: string;
  notice?: string;
  addMode?: number;
  status?: number;
}

export interface GroupListResp {
  code: number;
  message: string;
  data: GroupInfoResp[];
}

export interface MemberVo {
  userId: number;
  nickname: string;
  role: number;
  avatar?: string;
}

export interface GroupMemberListResp {
  code: number;
  message: string;
  memberList: MemberVo[];
}

export interface SearchGroupListResp {
  code: number;
  message: string;
  data: GroupInfoResp[];
}
