export type Message = {
  id: number;
  text: string;
  type: "user" | "system";
  roomId: number;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  createdAt: number;
};

export type MessageRaw = {
  id: number;
  content: string;
  type: "user" | "system";
  room_id: number;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  created_at: number;
};

export type Room = {
  name: string;
  id: number;
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
