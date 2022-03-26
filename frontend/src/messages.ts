import { Message, MessageRaw } from "./types/coreTypes";

export const parseMessage = (data: MessageRaw): Message => {
  return {
    id: data.id,
    type: data.type,
    text: data.content,
    roomId: data.room_id,
    createdAt: data.created_at,
    user: {
      id: data.user.id,
      name: data.user.name,
      avatar: data.user.avatar,
    },
    isInbound: data.is_inbound,
  };
};
