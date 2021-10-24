export type ConnectionState = "disconnected" | "connecting" | "connected";

export const connectWS = (
  onUpdateStatus: (status: ConnectionState) => void
) => {
  const connection = new WebSocket(import.meta.env.VITE_WS_DOMAIN + "/ws");
  onUpdateStatus("connecting");

  connection.onopen = () => {
    onUpdateStatus("connected");
  };

  connection.onclose = () => {
    const interval = setInterval(() => {
      onUpdateStatus("connecting");

      if (connection.readyState === WebSocket.OPEN) {
        clearInterval(interval);
      } else {
        connectWS(onUpdateStatus);
      }
    }, 10000);
  };

  connection.onerror = () => {
    onUpdateStatus("disconnected");
  };

  return connection;
};
