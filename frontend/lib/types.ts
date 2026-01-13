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

export interface AuthRequest {
  username: string;
  password: string;
}

export interface AuthResponse {
  access_token_expiration: number;
}

export interface CreateConversationRequest {
  conversation_name: string;
  conversation_id: string;
}

export interface StartDirectConversationRequest {
  to_user_id: string;
}

export interface StartDirectConversationResponse {
  conversation_id: string;
}

export interface ConversationIdRequest {
  conversation_id: string;
}

export interface InviteUserRequest {
  conversation_id: string;
  user_id: string;
}

export interface RenameConversationRequest {
  conversation_id: string;
  new_name: string;
}

export interface PaginationParams {
  page?: number;
  page_size?: number;
}

export interface WSIncomingMessage {
  type: "group_message" | "direct_message";
  data: {
    content: string;
    conversation_id: string;
  };
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
