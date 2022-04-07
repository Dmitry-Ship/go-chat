export type Message = {
  id: number;
  text: string;
  type: "user" | "system";
  conversationId: number;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  isInbound: boolean;
  createdAt: number;
};

export type MessageRaw = {
  id: number;
  text: string;
  type: "user" | "system";
  conversation_id: number;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  created_at: number;
  is_inbound: boolean;
};

export type Conversation = {
  name: string;
  id: number;
};

export type MessageEvent = {
  type: "message";
  data: MessageRaw;
};

export type User = {
  id: string;
  avatar: string;
  name: string;
};

export type Event = MessageEvent;
