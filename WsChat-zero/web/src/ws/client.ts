import { useChatStore } from '@/store/chat';
import type { MessageVo } from '@/types';

type WsMessageHandler = (msg: WsMessage) => void;

interface WsMessage {
  type: string;
  data: any;
}

class WsClient {
  private ws: WebSocket | null = null;
  private url: string = '';
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private handlers: Map<string, WsMessageHandler[]> = new Map();
  private shouldReconnect = true;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;

  connect(token: string) {
    this.shouldReconnect = true;
    this.reconnectAttempts = 0;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    // 开发环境直连代理
    this.url = `ws://localhost:8888/wss?token=${token}`;

    this.createConnection();
  }

  private createConnection() {
    if (this.ws) {
      this.ws.close();
    }

    try {
      this.ws = new WebSocket(this.url);
    } catch (e) {
      console.error('WS connect error:', e);
      this.scheduleReconnect();
      return;
    }

    this.ws.onopen = () => {
      console.log('WS connected');
      this.reconnectAttempts = 0;
      this.startHeartbeat();
    };

    this.ws.onmessage = (event) => {
      try {
        const msg: WsMessage = JSON.parse(event.data);
        this.dispatch(msg);
      } catch (e) {
        console.error('WS parse error:', e);
      }
    };

    this.ws.onclose = () => {
      console.log('WS disconnected');
      this.stopHeartbeat();
      if (this.shouldReconnect) {
        this.scheduleReconnect();
      }
    };

    this.ws.onerror = (e) => {
      console.error('WS error:', e);
    };
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('WS max reconnect attempts reached');
      return;
    }

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    this.reconnectAttempts++;

    this.reconnectTimer = setTimeout(() => {
      console.log(`WS reconnecting (attempt ${this.reconnectAttempts})...`);
      this.createConnection();
    }, delay);
  }

  private startHeartbeat() {
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

  sendMessage(content: string, receiverId: number, chatType: number, sessionId: number, msgType = 1) {
    return this.send({
      type: 'text',
      data: { content, receiver_id: receiverId, chat_type: chatType, msg_type: msgType, session_id: sessionId },
    });
  }

  private dispatch(msg: WsMessage) {
    const handlers = this.handlers.get(msg.type) || [];
    handlers.forEach((h) => h(msg));

    // "all" handler
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

  disconnect() {
    this.shouldReconnect = false;
    this.stopHeartbeat();
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

export const wsClient = new WsClient();
export default wsClient;
