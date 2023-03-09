import { getUser, login, logout, refreshToken, signup } from "../api/fetch";
import { User } from "../types/coreTypes";

export interface IAuthenticationService {
  login(username: string, password: string): Promise<void>;
  signup(username: string, password: string): Promise<void>;
  logout(): Promise<void>;
  onStateChanged(callback: (user: User | null, error: string) => void): void;
  rotateTokens(): void;
}

export class AuthenticationService implements IAuthenticationService {
  private onStateChangedCallback: (user: User | null, error: string) => void;
  private timeout!: NodeJS.Timeout;

  constructor() {
    this.onStateChangedCallback = (user: User | null, error: string) => {};
  }

  onStateChanged = (callback: (user: User | null, error: string) => void) => {
    this.onStateChangedCallback = callback;

    return () => clearTimeout(this.timeout);
  };

  logout = async () => {
    const result = await logout();

    if (result.status) {
      this.onStateChangedCallback(null, "");
      clearTimeout(this.timeout);
    } else {
      this.onStateChangedCallback(null, result.error || "Unknown error");
    }
  };

  login = async (username: string, password: string) => {
    const result = await login({
      username,
      password,
    });

    if (result.status) {
      this.fetchUser();
      this.rotateTokens();
    } else {
      this.onStateChangedCallback(null, result.error || "Unknown error");
    }
  };

  signup = async (username: string, password: string) => {
    const result = await signup({
      username,
      password,
    });

    if (result.status) {
      this.fetchUser();
      this.rotateTokens();
    } else {
      this.onStateChangedCallback(null, result.error || "Unknown error");
    }
  };

  rotateTokens = () => {
    const refresh = async () => {
      const result = await refreshToken();
      if (result.access_token_expiration) {
        const accessTokenRefreshInterval =
          result?.access_token_expiration / 1000000 / 2;

        this.fetchUser();

        setTimeout(refresh, accessTokenRefreshInterval);
      } else {
        clearTimeout(this.timeout);
        this.onStateChangedCallback(null, "");
      }
    };

    if (!this.timeout) {
      this.timeout = setTimeout(refresh, 0);
    }
  };

  private fetchUser = async () => {
    const getUserResult = await getUser()();

    if (getUserResult) {
      this.onStateChangedCallback(getUserResult, "");
    }
  };
}
