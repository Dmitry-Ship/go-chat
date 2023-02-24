"use client";
import axios from "axios";
import {
  Contact,
  ConversationFull,
  ConversationListItem,
  MessageRaw,
  User,
} from "../types/coreTypes";

export const makeCommand =
  <T>(url: string) =>
  async (body?: T): Promise<any> => {
    const result = await axios.post("/api" + url, body);

    return result.data;
  };

type conversationRequest = { conversation_id: string };

export const login = makeCommand<{
  username: string;
  password: string;
}>("/login");

export const logout = makeCommand("/logout");
export const signup = makeCommand<{
  username: string;
  password: string;
}>("/signup");

export const refreshToken = makeCommand("/refreshToken");

export const createConversation = makeCommand<
  {
    conversation_name: string;
  } & conversationRequest
>("/createConversation");

export const deleteConversation = makeCommand<conversationRequest>(
  "/deleteConversation"
);

export const renameConversation = makeCommand<
  {
    new_name: string;
  } & conversationRequest
>("/renameConversation");

export const leaveConversation =
  makeCommand<conversationRequest>("/leaveConversation");

export const startDirectConversation = makeCommand<{
  to_user_id: string;
}>("/startDirectConversation");

export const inviteUserToConversation = makeCommand<
  {
    user_id: string;
  } & conversationRequest
>("/inviteUserToConversation");

export const joinConversation =
  makeCommand<conversationRequest>("/joinConversation");

export const kick = makeCommand<{ user_id: string }>("/kick");

export const makeQuery =
  <T>(url: string) =>
  (param = "") =>
  async (): Promise<T> => {
    const { data } = await axios.get("/api" + url + param);
    console.log("data", data);

    return data as T;
  };

export const getConversation = makeQuery<ConversationFull>("/getConversation");
export const getPotentialInvitees = makeQuery<Contact[]>(
  "/getPotentialInvitees"
);
export const getParticipants = makeQuery<Contact[]>("/getParticipants");
export const getUser = makeQuery<User>("/getUser");

export const makePaginatedQuery = async <T>(
  url: string,
  page: number,
  pageSize = 50
): Promise<T[]> => {
  const paginationParams =
    (url.includes("?") ? "&" : "?") + "page=" + page + "&page_size=" + pageSize;

  const result = await makeQuery<T[]>(url)(paginationParams)();

  return result;
};

export const getContacts = (page: number, params: string = "") =>
  makePaginatedQuery<Contact>("/getContacts" + params, page, 50);

export const getConversations = (page: number, params: string = "") =>
  makePaginatedQuery<ConversationListItem>(
    "/getConversations" + params,
    page,
    50
  );

export const getConversationsMessages = (page: number, params: string = "") =>
  makePaginatedQuery<MessageRaw>(
    "/getConversationsMessages" + params,
    page,
    50
  );
