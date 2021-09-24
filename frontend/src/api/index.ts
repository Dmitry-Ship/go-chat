const connection = new WebSocket(import.meta.env.VITE_WS_DOMAIN + "/ws");

const connect = (cb: (msg: any) => void) => {
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

const sendMsg = <T>(msg: T) => {
  const stringifiedMessage = JSON.stringify(msg);
  connection.send(stringifiedMessage);
};

export { connect, sendMsg };
