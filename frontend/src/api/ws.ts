export type ConnectionState = "disconnected" | "connecting" | "connected";

export type IWSService = {
  send: (type: string, payload: Record<string, any>) => void;
  close: () => void;
  setOnUpdateStatus: (callback: (state: ConnectionState) => void) => void;
  onNotification: (event: string, callback: (data: any) => void) => void;
  connect: () => void;
  getStatus: () => ConnectionState;
};

export class WSService implements IWSService {
  private connection: WebSocket | null = null;
  private onUpdateStatus: (status: ConnectionState) => void = () => {};
  private status: ConnectionState = "disconnected";
  private events: Record<string, (data: any) => void> = {};

  public connect = () => {
    if (this.connection) {
      this.connection.close();
    }

    const connection = new WebSocket(process.env.NEXT_PUBLIC_WS_DOMAIN + "/ws");

    this.updateStatus("connecting");

    connection.onopen = () => {
      this.updateStatus("connected");
    };

    connection.onerror = () => {
      this.startReconnection();
      this.updateStatus("disconnected");
    };

    connection.onmessage = (event) => {
      const data = JSON.parse(event.data);

      this.events[data.type]?.(data);
    };

    this.connection = connection;
  };

  public startReconnection = () => {
    const intervalId = setInterval(() => {
      if (this.getStatus() === "disconnected") {
        this.connect();
      }

      if (this.getStatus() === "connected") {
        clearInterval(intervalId);
      }
    }, 5000);
  };

  private updateStatus = (status: ConnectionState) => {
    this.status = status;
    this.onUpdateStatus(status);
  };

  public setOnUpdateStatus = (cb: (status: ConnectionState) => void) => {
    this.onUpdateStatus = cb;
  };

  public send = (type: string, payload: Record<string, any>) => {
    const notificationObj = {
      type: type,
      data: payload,
    };

    const stringified = JSON.stringify(notificationObj);
    this.connection?.send(stringified);
  };

  public getStatus = () => this.status;

  public close = () => {
    this.connection?.close();
  };

  public onNotification = (event: string, callback: (data: any) => void) => {
    this.events[event] = callback;
  };
}
