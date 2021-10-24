import React, { createContext, useContext, useEffect, useState } from "react";
import { ConnectionState, connectWS } from "./api/ws";

type ws = {
  status: ConnectionState;
  sendNotification: (notification: {
    type: string;
    data: Record<string, any>;
  }) => void;
  subscribe: (event: string, cb: (msg: any) => void) => void;
};

export const useProvideWS = ({ isEnabled }: { isEnabled: boolean }): ws => {
  const [status, setStatus] = useState<ConnectionState>("disconnected");

  const [events, setEvents] = useState<{ [event: string]: (msg: any) => void }>(
    {}
  );

  const [connection, setConnection] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (isEnabled) {
      const conn = connectWS(setStatus);

      setConnection(conn);
    }
  }, [isEnabled]);

  if (connection !== null) {
    connection.onmessage = (msg) => {
      const data = JSON.parse(msg.data);

      events[data.type]?.(data);
    };
  }

  return {
    status,
    sendNotification: (notification) => {
      const stringifiedMessage = JSON.stringify(notification);
      connection?.send(stringifiedMessage);
    },
    subscribe: (event, cb) => {
      setEvents((prev) => ({ ...prev, [event]: cb }));
    },
  };
};

const wsContext = createContext<ws>({
  status: "disconnected",
  sendNotification: () => {},
  subscribe: () => {},
});

export const ProvideWS: React.FC<{
  children: React.ReactNode;
  isEnabled: boolean;
}> = ({ children, isEnabled }) => {
  const ws = useProvideWS({ isEnabled });
  return <wsContext.Provider value={ws}>{children}</wsContext.Provider>;
};

export const useWS = () => {
  return useContext(wsContext);
};
