let url = "/ws";
if (import.meta.env.DEV) {
  const { hostname, port } = window.location;
  url = `ws://${hostname}${port ? ":" + port : ""}` + url;
} else {
  url = import.meta.env.VITE_WS_DOMAIN + url;
}

const events: Record<string, (msg: any) => void> = {};

export const onEvent = (event: string, cb: (msg: any) => void) => {
  events[event] = cb;
};

export const connection = new WebSocket(url);

connection.onmessage = (event) => {
  const parsedMessage = JSON.parse(event.data);
  events[parsedMessage.type]?.(parsedMessage);
};

export const sendNotification = (notification: {
  type: string;
  data: Record<string, any>;
}) => {
  const stringifiedMessage = JSON.stringify(notification);
  connection.send(stringifiedMessage);
};

export const sendMsg = (msg: string, roomId: number, userId: number) => {
  sendNotification({
    type: "message",
    data: { content: msg, room_id: roomId, user_id: userId },
  });
};
