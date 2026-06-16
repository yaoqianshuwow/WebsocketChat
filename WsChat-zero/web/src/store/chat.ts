import { create } from 'zustand';
import api from '@/api/client';
import type { ContactVo, MessageVo, SessionResp, SessionVo } from '@/types';

type IncomingWsMessage = {
  sessionId?: number;
  senderId?: number;
  receiverId?: number;
  chatType?: number;
  msgType?: number;
  content?: string;
  fileUrl?: string;
  fileName?: string;
  fileSize?: number;
  sendName?: string;
  sendAvatar?: string;
  typing?: boolean;
};

let peerTypingTimer: ReturnType<typeof setTimeout> | null = null;

function sameMessageSeed(a: MessageVo, b: IncomingWsMessage) {
  return (
    a.senderId === (b.senderId || 0) &&
    a.receiverId === (b.receiverId || 0) &&
    a.msgType === (b.msgType || 0) &&
    (a.content || '') === (b.content || '') &&
    (a.fileUrl || '') === (b.fileUrl || '') &&
    (a.fileName || '') === (b.fileName || '') &&
    (a.fileSize || 0) === (b.fileSize || 0)
  );
}

interface ChatState {
  sessions: SessionVo[];
  currentSession: SessionVo | null;
  messages: MessageVo[];
  contacts: ContactVo[];
  loading: boolean;
  loadingHistory: boolean;
  messageBeforeId: number;
  messageHasMore: boolean;
  wsState: 'idle' | 'connecting' | 'connected' | 'reconnecting' | 'disconnected' | 'error';
  wsHint: string;
  peerTyping: boolean;

  loadSessions: () => Promise<void>;
  loadMessages: (sessionId: number) => Promise<void>;
  loadOlderMessages: () => Promise<void>;
  loadContacts: () => Promise<void>;
  setCurrentSession: (session: SessionVo | null) => void;
  addMessage: (msg: MessageVo) => void;
  syncIncomingMessage: (msg: IncomingWsMessage) => void;
  setPeerTyping: (sessionId: number | null, typing: boolean) => void;
  markMessageStatus: (localId: string, status: NonNullable<MessageVo['status']>) => void;
  setWsStatus: (state: ChatState['wsState'], hint?: string) => void;
  createSession: (peerId: number, sessionType: number, sessionName?: string) => Promise<boolean>;
}

export const useChatStore = create<ChatState>((set, get) => ({
  sessions: [],
  currentSession: null,
  messages: [],
  contacts: [],
  loading: false,
  loadingHistory: false,
  messageBeforeId: 0,
  messageHasMore: true,
  wsState: 'idle',
  wsHint: '未连接',
  peerTyping: false,

  loadSessions: async () => {
    try {
      const resp = await api.getSessionList();
      if (resp.code === 0) {
        set({ sessions: resp.data || [] });
      }
    } catch {
      // ignore
    }
  },

  loadMessages: async (sessionId) => {
    set({ loading: true });
    try {
      const session = get().sessions.find((item) => item.sessionId === sessionId) || get().currentSession;
      const size = 100;
      const resp =
        session?.sessionType === 2
          ? await api.getGroupMessageList(session.peerId, 0, size)
          : await api.getMessageList(sessionId, 0, size);
      if (resp.code === 0) {
        const msgs = resp.data || [];
        set({
          messages: msgs,
          messageBeforeId: msgs.length > 0 ? msgs[0].msgId || 0 : 0,
          messageHasMore: msgs.length >= size,
        });
      }
    } catch {
      // ignore
    }
    set({ loading: false });
  },

  loadOlderMessages: async () => {
    const session = get().currentSession;
    if (!session || get().loadingHistory || !get().messageHasMore) return;

    const beforeId = get().messageBeforeId;
    if (beforeId <= 0) {
      set({ messageHasMore: false });
      return;
    }
    const size = 20;
    set({ loadingHistory: true });
    try {
      const resp =
        session.sessionType === 2
          ? await api.getGroupMessageList(session.peerId, beforeId, size)
          : await api.getMessageList(session.sessionId, beforeId, size);
      if (resp.code === 0) {
        const incoming = resp.data || [];
        set((state) => {
          const seen = new Set(state.messages.map((item) => item.msgId || item.localId));
          const older = incoming.filter((item) => !seen.has(item.msgId || item.localId));
          return {
            messages: [...older, ...state.messages],
            messageBeforeId: older.length > 0 ? older[0].msgId || 0 : 0,
            messageHasMore: incoming.length >= size,
          };
        });
      }
    } catch {
      // ignore
    }
    set({ loadingHistory: false });
  },

  loadContacts: async () => {
    try {
      const resp = await api.getContactList();
      if (resp.code === 0) {
        set({ contacts: resp.data || [] });
      }
    } catch {
      // ignore
    }
  },

  setCurrentSession: (session) => {
    set({ currentSession: session, messages: [], messageBeforeId: 0, messageHasMore: true, peerTyping: false });
    if (session) {
      get().loadMessages(session.sessionId);
    }
  },

  addMessage: (msg) => {
    set((state) => ({ messages: [...state.messages, msg] }));
  },

  syncIncomingMessage: (msg) => {
    set((state) => {
      const session = state.currentSession;
      const incomingSessionId = msg.sessionId || 0;
      const matchesCurrentSession =
        !!session &&
        (
          (incomingSessionId > 0 && session.sessionId === incomingSessionId) ||
          (session.sessionType === 2 && msg.receiverId === session.peerId) ||
          (session.sessionType === 1 && (msg.senderId === session.peerId || msg.receiverId === session.peerId))
        );
      if (!matchesCurrentSession) {
        return state;
      }

      const normalized: MessageVo = {
        senderId: msg.senderId || 0,
        receiverId: msg.receiverId || 0,
        msgType: msg.msgType || 1,
        content: msg.content,
        fileUrl: msg.fileUrl,
        fileName: msg.fileName,
        fileSize: msg.fileSize,
        createdAt: Math.floor(Date.now() / 1000),
        status: 'sent',
        sendName: msg.sendName,
        sendAvatar: msg.sendAvatar,
        mine: msg.senderId === get().currentSession?.peerId ? false : undefined,
      };

      const nextMessages = [...state.messages];

      const localIndex = nextMessages.findIndex((item) => item.localId && sameMessageSeed(item, msg));
      if (localIndex >= 0) {
        const local = nextMessages[localIndex];
        nextMessages[localIndex] = {
          ...local,
          ...normalized,
          localId: local.localId,
          status: 'sent',
          mine: true,
        };
        return { messages: nextMessages };
      }

      if (nextMessages.some((item) => item.msgId && item.msgId === normalized.msgId)) {
        return state;
      }

      if (nextMessages.some((item) => sameMessageSeed(item, msg))) {
        return state;
      }

      return { messages: [...nextMessages, normalized] };
    });
  },

  setPeerTyping: (sessionId, typing) => {
    const currentSession = get().currentSession;
    if (!currentSession || currentSession.sessionId !== sessionId) {
      return;
    }
    if (peerTypingTimer) {
      clearTimeout(peerTypingTimer);
      peerTypingTimer = null;
    }

    if (!typing) {
      set({ peerTyping: false });
      return;
    }

    set({ peerTyping: true });
    peerTypingTimer = setTimeout(() => {
      const activeSession = get().currentSession;
      if (activeSession?.sessionId === sessionId) {
        set({ peerTyping: false });
      }
      peerTypingTimer = null;
    }, 1800);
  },

  markMessageStatus: (localId, status) => {
    set((state) => ({
      messages: state.messages.map((item) => (item.localId === localId ? { ...item, status } : item)),
    }));
  },

  setWsStatus: (wsState, hint) => {
    const defaultHintMap: Record<ChatState['wsState'], string> = {
      idle: '未连接',
      connecting: '正在连接',
      connected: '连接正常',
      reconnecting: '正在重连',
      disconnected: '连接已断开',
      error: '连接异常',
    };

    set({
      wsState,
      wsHint: hint || defaultHintMap[wsState],
    });
  },

  createSession: async (peerId, sessionType, sessionName) => {
    try {
      const existing = get().sessions.find((item) => item.peerId === peerId && item.sessionType === sessionType);
      if (existing) {
        const nextSession = {
          ...existing,
          sessionName: existing.sessionName || sessionName || existing.sessionName,
        };
        set({ currentSession: nextSession, messages: [] });
        await get().loadMessages(nextSession.sessionId);
        return true;
      }

      const resp: SessionResp = await api.createSession(peerId, sessionType, sessionName || '');
      if (resp.code === 0) {
        await get().loadSessions();
        const created =
          get().sessions.find((item) => item.peerId === peerId && item.sessionType === sessionType) ||
          (resp.sessionId
            ? {
                sessionId: resp.sessionId,
                peerId: resp.peerId || peerId,
                sessionType: resp.sessionType || sessionType,
                sessionName: resp.sessionName || sessionName || '',
                lastMsgContent: '',
                lastMsgTime: Math.floor(Date.now() / 1000),
                unreadCount: 0,
              }
            : null);
        if (created) {
          set({ currentSession: created, messages: [] });
          await get().loadMessages(created.sessionId);
        }
        return true;
      }
      return false;
    } catch {
      return false;
    }
  },
}));
