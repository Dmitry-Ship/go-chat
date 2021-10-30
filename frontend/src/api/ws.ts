export type ConnectionState = "disconnected" | "connecting" | "connected";

export const connectWS = (
  onUpdateStatus: (status: ConnectionState) => void
) => {
  const connection = new WebSocket(process.env.NEXT_PUBLIC_WS_DOMAIN + "/ws");

  onUpdateStatus("connecting");

  connection.onopen = () => {
    onUpdateStatus("connected");
  };

  connection.onerror = () => {
    onUpdateStatus("disconnected");
  };

  return connection;
};
