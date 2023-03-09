import React, { createContext, useContext, useState } from "react";
import { ConnectionState, WSService } from "../api/ws";

type ws = {
  status: ConnectionState;
  sendNotification: (type: string, payload: Record<string, any>) => void;
  onNotification: (event: string, cb: (msg: any) => void) => void;
};

const wsContext = createContext<ws | null>(null);

export const ProvideWS = ({ children }: { children: React.ReactNode }) => {
  const wsService = WSService.getInstance();

  const [status, setStatus] = useState<ConnectionState>(
    ConnectionState.CONNECTING
  );

  wsService.setOnUpdateStatus(setStatus);

  const value = {
    status,
    sendNotification: wsService.send,
    onNotification: wsService.onNotification,
  };

  return <wsContext.Provider value={value}>{children}</wsContext.Provider>;
};

export const useWebSocket = () => {
  return useContext(wsContext) as ws;
};
