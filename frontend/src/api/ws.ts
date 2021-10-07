let url = "/ws";
if (import.meta.env.DEV) {
  const { hostname, port } = window.location;
  url = `ws://${hostname}${port ? ":" + port : ""}` + url;
} else {
  url = import.meta.env.VITE_WS_DOMAIN + url;
}

const connection = new WebSocket(url);

export const connect = (cb: (msg: any) => void) => {
  connection.onopen = () => {
    console.log("Successfully Connected");
  };

  connection.onmessage = (event) => {
    const parsedMessage = JSON.parse(event.data);
    cb(parsedMessage);
  };

  connection.onclose = (event) => {
    console.log("Socket Closed Connection: ", event);
  };

  connection.onerror = (error) => {
    console.log("Socket Error: ", error);
  };
};

export const sendMsg = (msg: string, roomId: number) => {
  const stringifiedMessage = JSON.stringify({ content: msg, room_id: roomId });
  connection.send(stringifiedMessage);
};
