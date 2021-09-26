export type Message = {
  text: string;
  type: "user" | "system";
  sender: {
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
    sender: {
      id: string;
      avatar: string;
      name: string;
    };
    created_at: number;
  };
};

export type Event = MessageEvent;
