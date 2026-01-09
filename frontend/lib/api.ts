import {
  AuthRequest,
  AuthResponse,
  CreateConversationRequest,
  StartDirectConversationRequest,
  StartDirectConversationResponse,
  InviteUserRequest,
  RenameConversationRequest,
  UserDTO,
  ContactDTO,
  ConversationDTO,
  ConversationFullDTO,
  MessageDTO,
} from "./types";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

const fetchWithAuth = async <T = unknown>(url: string, options: RequestInit = {}): Promise<T> => {
  const response = await fetch(`${API_BASE}${url}`, {
    ...options,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.text().catch(() => response.statusText);
    throw new Error(error || response.statusText);
  }

  return response.json() as Promise<T>;
};

export const signup = (username: string, password: string) =>
  fetchWithAuth<AuthResponse>("/api/signup", {
    method: "POST",
    body: JSON.stringify({ username, password }),
  });

export const login = (username: string, password: string) =>
  fetchWithAuth<AuthResponse>("/api/login", {
    method: "POST",
    body: JSON.stringify({ username, password }),
  });

export const refreshToken = () =>
  fetchWithAuth<AuthResponse>("/api/refreshToken", { method: "POST" });

export const logout = () =>
  fetchWithAuth<string>("/api/logout", { method: "POST" });

export const getUser = () => fetchWithAuth<UserDTO>("/api/getUser");

export const getContacts = (page = 1, pageSize = 20) =>
  fetchWithAuth<ContactDTO[]>(
    `/api/getContacts?page=${page}&page_size=${pageSize}`
  );

export const createConversation = (name: string, id: string) =>
  fetchWithAuth<string>("/api/createConversation", {
    method: "POST",
    body: JSON.stringify({ conversation_name: name, conversation_id: id }),
  });

export const startDirectConversation = (toUserId: string) =>
  fetchWithAuth<StartDirectConversationResponse>("/api/startDirectConversation", {
    method: "POST",
    body: JSON.stringify({ to_user_id: toUserId }),
  });

export const getConversations = (page = 1, pageSize = 20) =>
  fetchWithAuth<ConversationDTO[]>(
    `/api/getConversations?page=${page}&page_size=${pageSize}`
  );

export const getConversation = (id: string) =>
  fetchWithAuth<ConversationFullDTO>(
    `/api/getConversation?conversation_id=${id}`
  );

export const getMessages = (
  conversationId: string,
  page = 1,
  pageSize = 20
) =>
  fetchWithAuth<MessageDTO[]>(
    `/api/getConversationsMessages?conversation_id=${conversationId}&page=${page}&page_size=${pageSize}`
  );

export const getParticipants = (
  conversationId: string,
  page = 1,
  pageSize = 20
) =>
  fetchWithAuth<UserDTO[]>(
    `/api/getParticipants?conversation_id=${conversationId}&page=${page}&page_size=${pageSize}`
  );

export const joinConversation = (conversationId: string) =>
  fetchWithAuth<string>("/api/joinConversation", {
    method: "POST",
    body: JSON.stringify({ conversation_id: conversationId }),
  });

export const leaveConversation = (conversationId: string) =>
  fetchWithAuth<string>("/api/leaveConversation", {
    method: "POST",
    body: JSON.stringify({ conversation_id: conversationId }),
  });

export const deleteConversation = (conversationId: string) =>
  fetchWithAuth<string>("/api/deleteConversation", {
    method: "POST",
    body: JSON.stringify({ conversation_id: conversationId }),
  });

export const renameConversation = (conversationId: string, newName: string) =>
  fetchWithAuth<string>("/api/renameConversation", {
    method: "POST",
    body: JSON.stringify({
      conversation_id: conversationId,
      new_name: newName,
    }),
  });

export const inviteUser = (conversationId: string, userId: string) =>
  fetchWithAuth<string>("/api/inviteUserToConversation", {
    method: "POST",
    body: JSON.stringify({ conversation_id: conversationId, user_id: userId }),
  });

export const kickUser = (conversationId: string, userId: string) =>
  fetchWithAuth<string>("/api/kick", {
    method: "POST",
    body: JSON.stringify({ conversation_id: conversationId, user_id: userId }),
  });

export const getPotentialInvitees = (
  conversationId: string,
  page = 1,
  pageSize = 20
) =>
  fetchWithAuth<UserDTO[]>(
    `/api/getPotentialInvitees?conversation_id=${conversationId}&page=${page}&page_size=${pageSize}`
  );
