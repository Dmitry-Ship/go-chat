// const connection = new WebSocket("ws://" + location.host + "/ws");
const connection = new WebSocket("ws://localhost:8080/ws");

const connect = (cb: (msg: any) => void) => {
  console.log("connecting");

  connection.onopen = () => {
    console.log("Successfully Connected");
  };

  connection.onmessage = (event) => {
    console.log("Message from WebSocket: ", event);
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
  console.log("sending msg: ", msg);
  const stringifiedMessage = JSON.stringify(msg);
  connection.send(stringifiedMessage);
};

export { connect, sendMsg };
