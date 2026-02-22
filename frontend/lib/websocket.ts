import { InfiniteData, QueryClient } from '@tanstack/react-query'
import {
  WSOutgoingMessage,
  MessageDTO,
  ConversationFullDTO,
  ConversationDTO,
  MessagePageResponse,
} from "./types";

type MessageHandler = (message: WSOutgoingMessage) => void;
type ErrorHandler = (error: Event) => void;
type CloseHandler = () => void;

class WebSocketManager {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;
  private reconnectDelay = 1000;
  private reconnectTimeoutId: ReturnType<typeof setTimeout> | null = null;
  private pendingDisconnectTimeoutId: ReturnType<typeof setTimeout> | null = null;
  private isManualClose = false;
  private messageHandlers: Set<MessageHandler> = new Set();
  private errorHandlers: Set<ErrorHandler> = new Set();
  private closeHandlers: Set<CloseHandler> = new Set();
  private url: string;
  private queryClient: QueryClient | null = null;

  constructor() {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:4000";
    const wsUrl = apiUrl.replace("http://", "ws://").replace("https://", "wss://");
    this.url = `${wsUrl}/ws`;
  }

  setQueryClient(client: QueryClient) {
    this.queryClient = client;
  }

  private clearReconnectTimeout() {
    if (this.reconnectTimeoutId) {
      clearTimeout(this.reconnectTimeoutId);
      this.reconnectTimeoutId = null;
    }
  }

  private clearPendingDisconnectTimeout() {
    if (this.pendingDisconnectTimeoutId) {
      clearTimeout(this.pendingDisconnectTimeoutId);
      this.pendingDisconnectTimeoutId = null;
    }
  }

  connect(): void {
    this.clearPendingDisconnectTimeout();

    if (this.ws?.readyState === WebSocket.OPEN || this.ws?.readyState === WebSocket.CONNECTING) {
      return;
    }

    this.isManualClose = false;
    this.clearReconnectTimeout();

    const socket = new WebSocket(this.url);
    this.ws = socket;

    socket.onopen = () => {
      console.log("WebSocket connected");
      this.reconnectAttempts = 0;
      this.reconnectDelay = 1000;
    };

    socket.onmessage = (event) => {
      try {
        const message: WSOutgoingMessage = JSON.parse(event.data);
        this.handleMessage(message);
        this.messageHandlers.forEach((handler) => handler(message));
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    socket.onerror = (error) => {
      if (this.isManualClose || this.ws !== socket) {
        return;
      }

      // Browser WebSocket errors have no actionable details; close event drives reconnection.
      console.warn("WebSocket transport error; waiting for reconnect.");
      this.errorHandlers.forEach((handler) => handler(error));
    };

    socket.onclose = (event) => {
      if (this.ws === socket) {
        this.ws = null;
      }

      console.log("WebSocket disconnected");
      this.closeHandlers.forEach((handler) => handler());

      if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
        console.log(`WebSocket closed (${event.code}). Reconnecting in ${this.reconnectDelay}ms...`);
        this.reconnectTimeoutId = setTimeout(() => {
          this.reconnectTimeoutId = null;
          this.connect();
        }, this.reconnectDelay);
      }
    };
  }

  private handleMessage(message: WSOutgoingMessage) {
    if (!this.queryClient) return;

    const events = message.events
      ? message.events
      : [{ type: message.type!, data: message.data! }];

    for (const event of events) {
      switch (event.type) {
        case 'message':
          const msg = event.data as MessageDTO;
          this.queryClient.setQueriesData(
            { queryKey: ['messages', msg.conversation_id] },
            (old: InfiniteData<MessagePageResponse> | undefined) => {
              if (!old) return old;
              const lastPageIndex = old.pages.length - 1;
              const lastPage = old.pages[lastPageIndex];
              if (lastPage.messages.some(m => m.id === msg.id)) return old;
              const updatedPages = [...old.pages];
              updatedPages[lastPageIndex] = {
                ...lastPage,
                messages: [...lastPage.messages, msg],
              };
              return { ...old, pages: updatedPages };
            }
          );
          this.queryClient.invalidateQueries({ queryKey: ['messages', msg.conversation_id] });
          this.queryClient.invalidateQueries({ queryKey: ['conversation-users', msg.conversation_id] });
          this.queryClient.invalidateQueries({ queryKey: ['conversations'] });
          break;

        case 'conversation_updated':
          const conv = event.data as ConversationFullDTO;
          this.queryClient.setQueryData(['conversation', conv.id], conv);
          this.queryClient.setQueryData(
            ['conversations', 1, 20],
            (old: ConversationDTO[] | undefined) => {
              if (!old) return old;
              return old.map(c =>
                c.id === conv.id
                  ? { ...c, name: conv.name, avatar: conv.avatar }
                  : c
              );
            }
          );
          this.queryClient.invalidateQueries({ queryKey: ['participants', conv.id] });
          break;

        case 'conversation_deleted':
          const { conversation_id } = event.data as { conversation_id: string };
          this.queryClient.removeQueries({ queryKey: ['conversation', conversation_id] });
          this.queryClient.removeQueries({ queryKey: ['messages', conversation_id] });
          this.queryClient.removeQueries({ queryKey: ['conversation-users', conversation_id] });
          this.queryClient.removeQueries({ queryKey: ['participants', conversation_id] });
          this.queryClient.invalidateQueries({ queryKey: ['conversations'] });
          break;
      }
    }
  }

  disconnect(): void {
    this.isManualClose = true;
    this.clearReconnectTimeout();
    this.clearPendingDisconnectTimeout();

    // Defer close to avoid development-mode remount churn closing a CONNECTING socket.
    this.pendingDisconnectTimeoutId = setTimeout(() => {
      this.pendingDisconnectTimeoutId = null;
      if (this.ws) {
        this.ws.close();
        this.ws = null;
      }
    }, 0);
  }

  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  onError(handler: ErrorHandler): () => void {
    this.errorHandlers.add(handler);
    return () => this.errorHandlers.delete(handler);
  }

  onClose(handler: CloseHandler): () => void {
    this.closeHandlers.add(handler);
    return () => this.closeHandlers.delete(handler);
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

export const wsManager = new WebSocketManager();
