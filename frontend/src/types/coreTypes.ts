export type Message = {
  text: string;
  type: "user" | "system";
  roomId: number;
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  created_at: number;
};

export type MessageEvent = {
  type: "message";
  data: {
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
};

export type AuthEvent = {
  type: "user_id";
  data: {
    user_id: string;
  };
};

export type Event = MessageEvent | AuthEvent;
