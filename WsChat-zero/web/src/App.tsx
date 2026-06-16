import { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import wsClient from '@/ws/client';
import Layout from '@/components/Layout';
import Login from '@/pages/Login';
import Register from '@/pages/Register';
import Chat from '@/pages/Chat';
import Contacts from '@/pages/Contacts';
import ContactSearch from '@/pages/ContactSearch';
import Groups from '@/pages/Groups';
import GroupSearch from '@/pages/GroupSearch';
import PptAssistant from '@/pages/PptAssistant';
import Profile from '@/pages/Profile';
import MessageSearch from '@/pages/MessageSearch';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((s) => s.token);
  if (!token && !localStorage.getItem('token')) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

export default function App() {
  const { init, token } = useAuthStore();

  useEffect(() => {
    init();
  }, []);

  useEffect(() => {
    if (token) {
      wsClient.connect(token);
      return () => { wsClient.disconnect(); };
    }
  }, [token]);

  useEffect(() => {
    const handleMessage = (msg: { data?: { sessionId?: number } }) => {
      const state = useChatStore.getState();
      state.loadSessions();
      state.syncIncomingMessage(msg?.data || {});
    };

    const handleTyping = (msg: { data?: { sessionId?: number; typing?: boolean } }) => {
      const state = useChatStore.getState();
      const data = msg?.data || {};
      state.setPeerTyping(data.sessionId || null, Boolean(data.typing));
    };

    wsClient.on('message:new', handleMessage);
    wsClient.on('typing', handleTyping);
    return () => {
      wsClient.off('message:new', handleMessage);
      wsClient.off('typing', handleTyping);
    };
  }, []);

  useEffect(() => {
    const handleAuthExpired = () => {
      useAuthStore.getState().logout();
      window.location.replace('/login');
    };
    window.addEventListener('wschat-auth-expired', handleAuthExpired);
    return () => window.removeEventListener('wschat-auth-expired', handleAuthExpired);
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route path="/chat" element={<Chat />} />
          <Route path="/contacts" element={<Contacts />} />
          <Route path="/contacts/search" element={<ContactSearch />} />
          <Route path="/groups" element={<Groups />} />
          <Route path="/groups/search" element={<GroupSearch />} />
          <Route path="/messages/search" element={<MessageSearch />} />
          <Route path="/ppt-assistant" element={<PptAssistant />} />
          <Route path="/profile" element={<Profile />} />
        </Route>
        <Route path="*" element={<Navigate to="/chat" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
