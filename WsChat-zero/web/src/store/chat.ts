import { create } from 'zustand';
import api from '@/api/client';
import type { SessionVo, MessageVo, ContactVo } from '@/types';

interface ChatState {
  sessions: SessionVo[];
  currentSession: SessionVo | null;
  messages: MessageVo[];
  contacts: ContactVo[];
  loading: boolean;

  loadSessions: () => Promise<void>;
  loadMessages: (sessionId: number) => Promise<void>;
  loadContacts: () => Promise<void>;
  setCurrentSession: (session: SessionVo | null) => void;
  addMessage: (msg: MessageVo) => void;
  createSession: (peerId: number, sessionType: number, sessionName?: string) => Promise<boolean>;
}

export const useChatStore = create<ChatState>((set, get) => ({
  sessions: [],
  currentSession: null,
  messages: [],
  contacts: [],
  loading: false,

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

  createSession: async (peerId, sessionType, sessionName) => {
    try {
      const resp = await api.createSession(peerId, sessionType, sessionName || '');
      if (resp.code === 0) {
        await get().loadSessions();
        return true;
      }
      return false;
    } catch {
      return false;
    }
  },
}));
