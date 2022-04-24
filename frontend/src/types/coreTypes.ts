export type BaseMessage = {
  id: string;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  conversationId: string;
  createdAt: string;
};

type JoinedMessage = BaseMessage & {
  type: "joined_conversation";
  text: string;
};

type LeftMessage = BaseMessage & {
  type: "left_conversation";
  text: string;
};

type RenamedMessage = BaseMessage & {
  type: "renamed_conversation";
  text: string;
  newName: string;
};

type TextMessage = BaseMessage & {
  text: string;
  type: "text";
  isInbound: boolean;
};

type baseMessageRaw = {
  id: string;
  conversation_id: string;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  created_at: string;
};

type JoinedMessageRaw = baseMessageRaw & {
  type: "joined_conversation";
};

type LeftMessageRaw = baseMessageRaw & {
  type: "left_conversation";
};

type TextMessageRaw = baseMessageRaw & {
  text: string;
  type: "text";
  is_inbound: boolean;
};

type RenamedMessageRaw = baseMessageRaw & {
  type: "renamed_conversation";
  new_name: string;
};

export type Message =
  | JoinedMessage
  | TextMessage
  | LeftMessage
  | RenamedMessage;

export type MessageRaw =
  | LeftMessageRaw
  | TextMessageRaw
  | JoinedMessageRaw
  | RenamedMessageRaw;

export type Conversation = {
  name: string;
  avatar: string;
  id: string;
  type: "private" | "public";
};

export type Contact = {
  name: string;
  avatar: string;
  id: string;
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
