export type BaseMessage = {
  id: string;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  conversationId: string;
  text: string;
  createdAt: string;
};

type JoinedMessage = BaseMessage & {
  type: "joined_conversation";
};

type InvitedMessage = BaseMessage & {
  type: "invited_conversation";
};

type LeftMessage = BaseMessage & {
  type: "left_conversation";
};

type RenamedMessage = BaseMessage & {
  type: "renamed_conversation";
};

type UnknownMessage = BaseMessage & {
  type: "unknown";
};

type TextMessage = BaseMessage & {
  type: "text";
  isInbound: boolean;
};

type baseMessageRaw = {
  id: string;
  conversation_id: string;
  text: string;
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

type InvitedMessageRaw = baseMessageRaw & {
  type: "invited_conversation";
};

type LeftMessageRaw = baseMessageRaw & {
  type: "left_conversation";
};

type TextMessageRaw = baseMessageRaw & {
  type: "text";
  is_inbound: boolean;
};

type RenamedMessageRaw = baseMessageRaw & {
  type: "renamed_conversation";
};

export type Message =
  | JoinedMessage
  | InvitedMessage
  | TextMessage
  | LeftMessage
  | RenamedMessage
  | UnknownMessage;

export type MessageRaw =
  | LeftMessageRaw
  | TextMessageRaw
  | JoinedMessageRaw
  | RenamedMessageRaw
  | InvitedMessageRaw;

export type Conversation = {
  name: string;
  avatar: string;
  id: string;
  type: "direct" | "group";
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
