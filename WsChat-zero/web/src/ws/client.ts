import { useChatStore } from '@/store/chat';

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
    useChatStore.getState().setWsStatus('connecting');

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    // 开发环境直连代理
    this.url = `${protocol}//${host}/wss?token=${token}`;

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
      useChatStore.getState().setWsStatus('error', '连接建立失败');
      this.scheduleReconnect();
      return;
    }

    this.ws.onopen = () => {
      console.log('WS connected');
      this.reconnectAttempts = 0;
      useChatStore.getState().setWsStatus('connected');
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
       useChatStore.getState().setWsStatus(this.shouldReconnect ? 'reconnecting' : 'disconnected');
      if (this.shouldReconnect) {
        this.scheduleReconnect();
      }
    };

    this.ws.onerror = (e) => {
      console.error('WS error:', e);
      useChatStore.getState().setWsStatus('error');
    };
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('WS max reconnect attempts reached');
      useChatStore.getState().setWsStatus('disconnected', '重连失败');
      return;
    }

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    this.reconnectAttempts++;
    useChatStore.getState().setWsStatus('reconnecting', `正在重连（第 ${this.reconnectAttempts} 次）`);

    this.reconnectTimer = setTimeout(() => {
      console.log(`WS reconnecting (attempt ${this.reconnectAttempts})...`);
      useChatStore.getState().setWsStatus('connecting', '重新建立连接');
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
    useChatStore.getState().setWsStatus('disconnected');
  }
}

export const wsClient = new WsClient();
export default wsClient;
