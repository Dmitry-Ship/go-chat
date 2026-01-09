import { QueryClient } from '@tanstack/react-query'
import {
  WSIncomingMessage,
  WSOutgoingMessage,
  MessageDTO,
  ConversationFullDTO,
  ConversationDTO,
} from "./types";

type MessageHandler = (message: WSOutgoingMessage) => void;
type ErrorHandler = (error: Event) => void;
type CloseHandler = () => void;

class WebSocketManager {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;
  private reconnectDelay = 1000;
  private isManualClose = false;
  private messageHandlers: Set<MessageHandler> = new Set();
  private errorHandlers: Set<ErrorHandler> = new Set();
  private closeHandlers: Set<CloseHandler> = new Set();
  private url: string;
  private queryClient: QueryClient | null = null;

  constructor() {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const wsUrl = apiUrl.replace("http://", "ws://").replace("https://", "wss://");
    this.url = `${wsUrl}/ws`;
  }

  setQueryClient(client: QueryClient) {
    this.queryClient = client;
  }

  connect(): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return;
    }

    this.isManualClose = false;
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log("WebSocket connected");
      this.reconnectAttempts = 0;
      this.reconnectDelay = 1000;
    };

    this.ws.onmessage = (event) => {
      try {
        const message: WSOutgoingMessage = JSON.parse(event.data);
        this.handleMessage(message);
        this.messageHandlers.forEach((handler) => handler(message));
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      this.errorHandlers.forEach((handler) => handler(error));
    };

    this.ws.onclose = () => {
      console.log("WebSocket disconnected");
      this.closeHandlers.forEach((handler) => handler());

      if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
        console.log(`Reconnecting in ${this.reconnectDelay}ms...`);
        setTimeout(() => this.connect(), this.reconnectDelay);
      }
    };
  }

  private handleMessage(message: WSOutgoingMessage) {
    if (!this.queryClient) return;

    switch (message.type) {
      case 'message':
        const msg = message.data as MessageDTO;
        this.queryClient.setQueryData(
          ['messages', msg.conversation_id, 1, 20],
          (old: MessageDTO[] | undefined) => {
            if (!old) return [msg];
            if (old.some(m => m.id === msg.id)) return old;
            return [...old, msg];
          }
        );
        this.queryClient.invalidateQueries({ queryKey: ['conversations'] });
        break;

      case 'conversation_updated':
        const conv = message.data as ConversationFullDTO;
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
        const { conversation_id } = message.data as { conversation_id: string };
        this.queryClient.removeQueries({ queryKey: ['conversation', conversation_id] });
        this.queryClient.removeQueries({ queryKey: ['messages', conversation_id] });
        this.queryClient.removeQueries({ queryKey: ['participants', conversation_id] });
        this.queryClient.invalidateQueries({ queryKey: ['conversations'] });
        break;
    }
  }

  disconnect(): void {
    this.isManualClose = true;
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  sendMessage(type: "group_message" | "direct_message", data: { content: string; conversation_id: string }): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error("WebSocket is not connected");
      return;
    }

    const message: WSIncomingMessage = { type, data };
    this.ws.send(JSON.stringify(message));
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
