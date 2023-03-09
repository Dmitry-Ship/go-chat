import { Message, MessageRaw, BaseMessage } from "./types/coreTypes";

function ParseBaseMessage(raw: MessageRaw): BaseMessage {
  return {
    id: raw.id,
    conversationId: raw.conversation_id,
    user: {
      id: raw.user.id,
      avatar: raw.user.avatar,
      name: raw.user.name,
    },
    text: raw.text,
    createdAt: raw.created_at,
  };
}

export function parseMessage(data: MessageRaw): Message {
  const base = ParseBaseMessage(data);
  switch (data.type) {
    case "joined_conversation":
      return {
        ...base,
        type: "joined_conversation",
      };
    case "left_conversation":
      return {
        ...base,
        type: "left_conversation",
      };

    case "renamed_conversation":
      return {
        ...base,
        type: "renamed_conversation",
      };

    case "text":
      return {
        ...base,
        type: "text",
        isInbound: data.is_inbound,
      };
    case "invited_conversation":
      return {
        ...base,
        type: "invited_conversation",
      };
    default:
      return {
        ...base,
        type: "unknown",
      };
  }
}
