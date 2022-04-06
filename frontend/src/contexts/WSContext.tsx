import React, { createContext, useContext, useState } from "react";
import { ConnectionState, IWSService, WSService } from "../api/ws";

type ws = {
  status: ConnectionState;
  sendNotification: (type: string, payload: Record<string, any>) => void;
  onNotification: (event: string, cb: (msg: any) => void) => void;
};

export const useProvideWS = (wsService: IWSService): ws => {
  const [status, setStatus] = useState<ConnectionState>(
    ConnectionState.CONNECTING
  );

  wsService.setOnUpdateStatus(setStatus);

  return {
    status,
    sendNotification: wsService.send,
    onNotification: wsService.onNotification,
  };
};

const wsContext = createContext<ws>({
  status: 0,
  sendNotification: () => {},
  onNotification: () => {},
});

export const ProvideWS: React.FC = ({ children }) => {
  const wsService = WSService.getInstance();

  const ws = useProvideWS(wsService);
  return <wsContext.Provider value={ws}>{children}</wsContext.Provider>;
};

export const useWS = () => {
  return useContext(wsContext);
};
