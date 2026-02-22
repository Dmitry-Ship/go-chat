export type MessageType = "user" | "system";

export interface UserDTO {
  id: string;
  avatar: string;
  name: string;
}

export interface MessageDTO {
  id: string;
  created_at: string;
  text: string;
  type: MessageType;
  user_id: string;
  conversation_id: string;
}

export interface MessagePageResponse {
  messages: MessageDTO[];
  next_cursor?: string;
  has_more: boolean;
}

export interface ConversationUsersResponse {
  users: Record<string, UserDTO>;
}

export interface ConversationDTO {
  id: string;
  name: string;
  avatar: string;
  type: "group" | "direct";
  last_message: MessageDTO | null;
}

export interface ConversationPageResponse {
  conversations: ConversationDTO[];
  next_cursor?: string;
  has_more: boolean;
}

export interface ConversationFullDTO {
  id: string;
  name: string;
  avatar: string;
  created_at: string;
  type: "group" | "direct";
  joined: boolean;
  participants_count: number;
  is_owner: boolean;
}

export interface ContactDTO {
  id: string;
  avatar: string;
  name: string;
}

export interface AuthResponse {
  access_token_expiration: number;
}

export interface StartDirectConversationResponse {
  conversation_id: string;
}

export interface WSNotificationEvent {
  type: "message" | "conversation_updated" | "conversation_deleted";
  data: MessageDTO | ConversationFullDTO | { conversation_id: string };
}

export interface WSOutgoingMessage {
  user_id: string;
  type?: "message" | "conversation_updated" | "conversation_deleted";
  data?: MessageDTO | ConversationFullDTO | { conversation_id: string };
  events?: WSNotificationEvent[];
}
