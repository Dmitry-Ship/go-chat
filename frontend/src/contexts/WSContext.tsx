import React, { createContext, useContext, useEffect, useState } from "react";
import { ConnectionState, IWSService, WSService } from "../api/ws";

type ws = {
  status: ConnectionState;
  sendNotification: (type: string, payload: Record<string, any>) => void;
  onNotification: (event: string, cb: (msg: any) => void) => void;
};

export const useProvideWS = (wsService: IWSService): ws => {
  const [status, setStatus] = useState<ConnectionState>(wsService.getStatus());

  wsService.setOnUpdateStatus(setStatus);

  useEffect(() => {
    wsService.connect();
    return () => {
      wsService.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return {
    status,
    sendNotification: wsService.send,
    onNotification: wsService.onNotification,
  };
};

const wsContext = createContext<ws>({
  status: "disconnected",
  sendNotification: () => {},
  onNotification: () => {},
});

export const ProvideWS: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const wsService = new WSService();
  const ws = useProvideWS(wsService);
  return <wsContext.Provider value={ws}>{children}</wsContext.Provider>;
};

export const useWS = () => {
  return useContext(wsContext);
};
