import { create } from 'zustand';
import api from '@/api/client';
import type { SessionVo, MessageVo, ContactVo, SessionResp } from '@/types';

interface ChatState {
  sessions: SessionVo[];
  currentSession: SessionVo | null;
  messages: MessageVo[];
  contacts: ContactVo[];
  loading: boolean;
  wsState: 'idle' | 'connecting' | 'connected' | 'reconnecting' | 'disconnected' | 'error';
  wsHint: string;

  loadSessions: () => Promise<void>;
  loadMessages: (sessionId: number) => Promise<void>;
  loadContacts: () => Promise<void>;
  setCurrentSession: (session: SessionVo | null) => void;
  addMessage: (msg: MessageVo) => void;
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
  wsState: 'idle',
  wsHint: '未连接',

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
      const resp = await api.getMessageList(sessionId);
      if (resp.code === 0) {
        set({ messages: resp.data || [] });
      }
    } catch {
      // ignore
    }
    set({ loading: false });
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
    set({ currentSession: session, messages: [] });
    if (session) {
      get().loadMessages(session.sessionId);
    }
  },

  addMessage: (msg) => {
    set((state) => ({ messages: [...state.messages, msg] }));
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
        const created = get().sessions.find((item) => item.peerId === peerId && item.sessionType === sessionType) || (resp.sessionId
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
