import { create } from 'zustand';
import api from '@/api/client';
import wsClient from '@/ws/client';
import type { LoginResp, UserInfoResp } from '@/types';

interface AuthState {
  token: string | null;
  user: (LoginResp & { userId?: number }) | null;
  userInfo: UserInfoResp | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<boolean>;
  register: (username: string, password: string, nickname?: string) => Promise<boolean>;
  logout: () => void;
  loadUserInfo: () => Promise<void>;
  init: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: null,
  user: null,
  userInfo: null,
  loading: false,

  init: () => {
    const token = localStorage.getItem('token');
    if (token) {
      api.setToken(token);
      set({ token });
      get().loadUserInfo();
    }
  },

  login: async (username, password) => {
    set({ loading: true });
    try {
      const resp = await api.login({ username, password });
      if (resp.code === 0 && resp.token) {
        api.setToken(resp.token);
        set({ token: resp.token, user: resp, loading: false });
        get().loadUserInfo();
        return true;
      }
      set({ loading: false });
      return false;
    } catch {
      set({ loading: false });
      return false;
    }
  },

  register: async (username, password, nickname) => {
    set({ loading: true });
    try {
      const resp = await api.register({ username, password, nickname });
      set({ loading: false });
      return resp.code === 0;
    } catch {
      set({ loading: false });
      return false;
    }
  },

  logout: () => {
    wsClient.disconnect();
    api.setToken(null);
    set({ token: null, user: null, userInfo: null });
  },

  loadUserInfo: async () => {
    try {
      const resp = await api.getUserInfo();
      if (resp.code === 0) {
        set({ userInfo: resp });
      }
    } catch {
      // ignore
    }
  },
}));
