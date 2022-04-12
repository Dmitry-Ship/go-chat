import { Message, MessageRaw, BaseMessage } from "./types/coreTypes";

const ParseBaseMessage = (raw: MessageRaw): BaseMessage => {
  return {
    id: raw.id,
    conversationId: raw.conversation_id,
    user: {
      id: raw.user.id,
      avatar: raw.user.avatar,
      name: raw.user.name,
    },
    createdAt: raw.created_at,
  };
};

export const parseMessage = (data: MessageRaw): Message => {
  const base = ParseBaseMessage(data);
  switch (data.type) {
    case "joined_conversation":
      return {
        ...base,
        text: `${base.user.name} joined the conversation`,
        type: "joined_conversation",
      };
    case "left_conversation":
      return {
        ...base,
        text: `${base.user.name} left the conversation`,
        type: "left_conversation",
      };

    case "renamed_conversation":
      return {
        ...base,
        text: `${base.user.name} renamed the conversation to ${data.new_name}`,
        type: "renamed_conversation",
        newName: data.new_name,
      };

    case "text":
      return {
        ...base,
        text: data.text,
        type: "text",
        isInbound: data.is_inbound,
      };

    default:
      throw new Error("Unknown message type");
  }
};
