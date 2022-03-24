import { makeCommand, makeQuery } from "../api/fetch";

export const rotateTokens = (cb: (isAuthenticated: boolean) => void) => {
  const refresh = async () => {
    const result = await makeCommand("/refreshToken");

    const accessTokenRefreshInterval =
      result.data?.access_token_expiration / 1000000 / 2;

    if (result.status) {
      setTimeout(refresh, accessTokenRefreshInterval);
    }

    cb(result.status);
  };

  const timeoutId = setTimeout(refresh, 0);

  return timeoutId;
};

export const fetchUser = async () => {
  const getUserResult = await makeQuery("/getUser");

  if (getUserResult.status) {
    return getUserResult.data;
  }

  return null;
};

export const login = async (
  username: string,
  password: string,
  cb: () => void
) => {
  const result = await makeCommand("/login", {
    username,
    password,
  });

  if (result.status) {
    cb();
  }
};

export const logout = async (cb: () => void) => {
  const result = await makeCommand("/logout");

  if (result.status) {
    cb();
  }
};

export const signup = async (
  username: string,
  password: string,
  cb: () => void
) => {
  const result = await makeCommand("/signup", {
    username,
    password,
  });

  if (result.status) {
    cb();
  }
};
