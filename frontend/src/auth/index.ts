import { makeCommand, makeQuery } from "../api/fetch";
import { User } from "../types/coreTypes";

export interface IAuthenticationService {
  login(username: string, password: string): Promise<void>;
  signup(username: string, password: string): Promise<void>;
  logout(): Promise<void>;
  onLogin(callback: () => void): void;
  onLogout(callback: () => void): void;
  fetchUser(): Promise<User>;
  rotateTokens(): NodeJS.Timeout;
}
export class AuthenticationService implements IAuthenticationService {
  private onLoginCallback: () => void;
  private onLogoutCallback: () => void;

  constructor() {
    this.onLoginCallback = () => {};
    this.onLogoutCallback = () => {};
  }

  onLogin = (callback: () => void) => {
    this.onLoginCallback = callback;
  };

  onLogout = (callback: () => void) => {
    this.onLogoutCallback = callback;
  };

  logout = async () => {
    const result = await makeCommand("/logout");

    if (result.status) {
      this.onLogoutCallback();
    }
  };

  login = async (username: string, password: string) => {
    const result = await makeCommand("/login", {
      username,
      password,
    });

    if (result.status) {
      this.onLoginCallback();
    }
  };

  signup = async (username: string, password: string) => {
    const result = await makeCommand("/signup", {
      username,
      password,
    });

    if (result.status) {
      this.onLoginCallback();
    }
  };

  rotateTokens = () => {
    const refresh = async () => {
      const result = await makeCommand("/refreshToken");

      const accessTokenRefreshInterval =
        result.data?.access_token_expiration / 1000000 / 2;

      if (result.status) {
        this.onLoginCallback();
        setTimeout(refresh, accessTokenRefreshInterval);
      } else {
        this.onLogoutCallback();
      }
    };

    const timeoutId = setTimeout(refresh, 0);

    return timeoutId;
  };

  fetchUser = async () => {
    const getUserResult = await makeQuery("/getUser");

    if (getUserResult.status) {
      return getUserResult.data;
    }

    return null;
  };
}
