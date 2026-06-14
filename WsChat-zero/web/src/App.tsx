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
import Groups from '@/pages/Groups';
import Profile from '@/pages/Profile';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((s) => s.token);
  if (!token && !localStorage.getItem('token')) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

export default function App() {
  const { init, token } = useAuthStore();
  const { loadSessions, loadMessages, currentSession } = useChatStore();

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
      void loadSessions();
      const sessionId = Number(msg?.data?.sessionId || 0);
      if (sessionId > 0 && currentSession?.sessionId === sessionId) {
        void loadMessages(sessionId);
      }
    };

    wsClient.on('message:new', handleMessage);
    return () => {
      wsClient.off('message:new', handleMessage);
    };
  }, [currentSession, loadMessages, loadSessions]);

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
          <Route path="/groups" element={<Groups />} />
          <Route path="/profile" element={<Profile />} />
        </Route>
        <Route path="*" element={<Navigate to="/chat" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
