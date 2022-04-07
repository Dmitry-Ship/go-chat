type SystemMessage = {
  id: string;
  type: 1;
  text: string;
  conversationId: string;
  createdAt: string;
};

type UserMessage = {
  id: string;
  text: string;
  type: 0;
  conversationId: string;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  isInbound: boolean;
  createdAt: string;
};

type SystemMessageRaw = {
  id: string;
  text: string;
  type: 1;
  conversation_id: string;
  created_at: string;
};

type UserMessageRaw = {
  id: string;
  text: string;
  type: 0;
  conversation_id: string;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  created_at: string;
  is_inbound: boolean;
};

export type Message = SystemMessage | UserMessage;

export type MessageRaw = SystemMessageRaw | UserMessageRaw;

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
