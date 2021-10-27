export type ConnectionState = "disconnected" | "connecting" | "connected";

export const connectWS = (
  onUpdateStatus: (status: ConnectionState) => void
) => {
  const connection = new WebSocket(import.meta.env.VITE_WS_DOMAIN + "/ws");

  onUpdateStatus("connecting");

  connection.onopen = () => {
    onUpdateStatus("connected");
  };

  connection.onerror = () => {
    onUpdateStatus("disconnected");
  };

  return connection;
};
