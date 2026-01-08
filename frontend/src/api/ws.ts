export enum ConnectionState {
  CONNECTING = 0,
  OPEN = 1,
  CLOSING = 2,
  CLOSED = 3,
}

export type IWSService = {
  send: (type: string, payload: Record<string, any>) => void;
  setOnUpdateStatus: (callback: (state: ConnectionState) => void) => void;
  onNotification: (event: string, callback: (data: any) => void) => void;
};

export class WSService implements IWSService {
  private connection: WebSocket | null = null;
  private onUpdateStatus: (status: ConnectionState) => void = () => {};
  private events: Record<string, (data: any) => void> = {};

  private static instance: IWSService;

  public static getInstance(): IWSService {
    if (!WSService.instance) {
      WSService.instance = new WSService();
    }

    return WSService.instance;
  }

  constructor() {
    this.connect();
  }

  private connect = () => {
    if (this.connection?.readyState === ConnectionState.OPEN) {
      this.connection.close();
      this.connection = null;
    }

    const connection = new WebSocket("ws:" + location.host + "/ws");

    this.updateStatus();

    connection.onopen = () => this.updateStatus();

    connection.onerror = () => {
      this.updateStatus();
      this.startReconnection();
    };

    connection.onmessage = (event) => {
      const data = JSON.parse(event.data);

      this.events[data.type]?.(data);
    };

    this.connection = connection;
  };

  private startReconnection = () => {
    const intervalId = setInterval(() => {
      if (this.connection?.readyState === ConnectionState.CLOSED) {
        this.connect();
      }

      if (this.connection?.readyState === ConnectionState.OPEN) {
        clearInterval(intervalId);
      }
    }, RECONNECTION_INTERVAL_MS);
  };

  private updateStatus = () => {
    this.onUpdateStatus(
      (this.connection?.readyState ??
        ConnectionState.CONNECTING) as ConnectionState
    );
  };

  public setOnUpdateStatus = (cb: (status: ConnectionState) => void) => {
    this.onUpdateStatus = cb;
  };

  public send = (type: string, payload: Record<string, any>) => {
    const data = JSON.stringify({ type, data: payload });
    this.connection?.send(data);
  };

  public onNotification = (event: string, callback: (data: any) => void) => {
    this.events[event] = callback;
  };
}
