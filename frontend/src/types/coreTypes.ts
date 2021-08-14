export type Message = {
  text: string;
  type: "user" | "system";
  sender: {
    id: string;
    avatar: string;
    name: string;
  };
  created_at: number;
  avatar: string;
};
