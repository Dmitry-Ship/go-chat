const connection = new WebSocket(import.meta.env.VITE_WS_DOMAIN + "/ws");

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

export const sendMsg = (msg: string) => {
  const stringifiedMessage = JSON.stringify({ content: msg });
  connection.send(stringifiedMessage);
};
