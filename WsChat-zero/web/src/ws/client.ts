import { useChatStore } from '@/store/chat';

type WsMessageHandler = (msg: WsMessage) => void;

interface WsMessage {
  type: string;
  data: any;
}

class WsClient {
  private ws: WebSocket | null = null;
  private url = '';
  private activeToken: string | null = null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private handlers: Map<string, WsMessageHandler[]> = new Map();
  private shouldReconnect = true;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;
  private connectionSeq = 0;

  connect(token: string) {
    if (this.activeToken === token && this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) {
      return;
    }

    this.disconnect(false);
    this.activeToken = token;
    this.shouldReconnect = true;
    this.reconnectAttempts = 0;
    useChatStore.getState().setWsStatus('connecting');

    const defaultProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const defaultHost = window.location.host;
    const baseUrl = import.meta.env.VITE_WS_BASE_URL || `${defaultProtocol}//${defaultHost}`;
    this.url = `${baseUrl}/wss?token=${token}`;

    this.createConnection();
  }

  private createConnection() {
    const seq = ++this.connectionSeq;
    const ws = new WebSocket(this.url);
    this.ws = ws;

    ws.onopen = () => {
      if (this.connectionSeq !== seq || this.ws !== ws) return;
      this.reconnectAttempts = 0;
      useChatStore.getState().setWsStatus('connected');
      this.startHeartbeat();
    };

    ws.onmessage = (event) => {
      if (this.connectionSeq !== seq || this.ws !== ws) return;
      try {
        const msg: WsMessage = JSON.parse(event.data);
        this.dispatch(msg);
      } catch (e) {
        console.error('WS parse error:', e);
      }
    };

    ws.onclose = () => {
      if (this.connectionSeq !== seq || this.ws !== ws) return;
      this.stopHeartbeat();
      this.ws = null;
      if (this.shouldReconnect) {
        useChatStore.getState().setWsStatus('reconnecting');
        this.scheduleReconnect();
      } else {
        useChatStore.getState().setWsStatus('disconnected');
      }
    };

    ws.onerror = (e) => {
      if (this.connectionSeq !== seq || this.ws !== ws) return;
      console.error('WS error:', e);
      useChatStore.getState().setWsStatus('error');
    };
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      useChatStore.getState().setWsStatus('disconnected', '重连失败');
      return;
    }

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    this.reconnectAttempts++;
    useChatStore.getState().setWsStatus('reconnecting', `正在重连（第 ${this.reconnectAttempts} 次）`);

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
    }

    this.reconnectTimer = setTimeout(() => {
      useChatStore.getState().setWsStatus('connecting', '重新建立连接');
      this.createConnection();
    }, delay);
  }

  private startHeartbeat() {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      this.send({ type: 'heartbeat', data: {} });
    }, 30000);
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  send(msg: WsMessage) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg));
      return true;
    }
    return false;
  }

  sendMessage(
    content: string,
    receiverId: number,
    chatType: number,
    sessionId: number,
    msgType = 1,
    extra?: { fileUrl?: string; fileName?: string; fileSize?: number },
  ) {
    return this.send({
      type: 'text',
      data: {
        content,
        receiver_id: receiverId,
        chat_type: chatType,
        msg_type: msgType,
        session_id: sessionId,
        file_url: extra?.fileUrl,
        file_name: extra?.fileName,
        file_size: extra?.fileSize,
      },
    });
  }

  private dispatch(msg: WsMessage) {
    const handlers = this.handlers.get(msg.type) || [];
    handlers.forEach((h) => h(msg));

    const allHandlers = this.handlers.get('*') || [];
    allHandlers.forEach((h) => h(msg));
  }

  on(type: string, handler: WsMessageHandler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, []);
    }
    this.handlers.get(type)!.push(handler);
  }

  off(type: string, handler: WsMessageHandler) {
    const handlers = this.handlers.get(type);
    if (handlers) {
      const idx = handlers.indexOf(handler);
      if (idx >= 0) handlers.splice(idx, 1);
    }
  }

  disconnect(clearStatus = true) {
    this.shouldReconnect = false;
    this.activeToken = null;
    this.connectionSeq++;
    this.stopHeartbeat();
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      const ws = this.ws;
      this.ws = null;
      try {
        ws.close();
      } catch {
        // ignore
      }
    }
    if (clearStatus) {
      useChatStore.getState().setWsStatus('disconnected');
    }
  }
}

export const wsClient = new WsClient();
export default wsClient;
