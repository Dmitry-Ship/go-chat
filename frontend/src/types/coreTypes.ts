export type Message = {
  text: string;
  type: "user" | "system";
  sender: number;
};
