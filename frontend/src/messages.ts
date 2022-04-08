import { Message, MessageRaw } from "./types/coreTypes";

export const parseMessage = (data: MessageRaw): Message => {
  if (data.type === "user") {
    return {
      id: data.id,
      type: data.type,
      text: data.text,
      conversationId: data.conversation_id,
      createdAt: data.created_at,
      user: {
        id: data.user.id,
        name: data.user.name,
        avatar: data.user.avatar,
      },
      isInbound: data.is_inbound,
    };
  }

  return {
    id: data.id,
    type: data.type,
    text: data.text,
    conversationId: data.conversation_id,
    createdAt: data.created_at,
  };
};
