import React, { createContext, useContext, useEffect, useState } from "react";
import { makeCommand, makeQuery } from "../api/fetch";
import { User } from "../types/coreTypes";

type auth = {
  user: null | User;
  isAuthenticated: boolean;
  isChecking: boolean;
  signup: (username: string, password: string) => void;
  login: (username: string, password: string) => void;
  logout: () => void;
};

export const useProvideAuth = (): auth => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isChecking, setIsChecking] = useState<boolean>(true);
  const [accessTokenRefreshInterval, setAccessTokenRefreshInterval] =
    useState<Number>(0);
  const [user, setUser] = useState<User | null>(null);

  const refreshToken = async () => {
    const result = await makeCommand("/refreshToken");

    const accessTokenRefreshInterval =
      result.data?.access_token_expiration / 1000000 / 2;

    setAccessTokenRefreshInterval(accessTokenRefreshInterval);
    setIsAuthenticated(Boolean(result.status));
  };

  useEffect(async () => {
    if (!isAuthenticated) {
      await refreshToken();
    }

    setIsChecking(false);

    if (isAuthenticated) {
      const getUserResult = await makeQuery("/getUser");

      if (getUserResult.status) {
        setUser(getUserResult.data);
      }

      const interval = setInterval(refreshToken, accessTokenRefreshInterval);

      return () => clearInterval(interval);
    }
  }, [isAuthenticated, accessTokenRefreshInterval]);

  return {
    isAuthenticated,
    isChecking,
    user,
    login: async (username: string, password: string) => {
      const result = await makeCommand("/login", {
        username: username,
        password,
      });

      if (result.status) {
        window.location.reload();
      }
    },
    logout: async () => {
      const result = await makeCommand("/logout");

      if (result.status) {
        window.location.reload();
      }
    },
    signup: async (username: string, password: string) => {
      const result = await makeCommand("/signup", {
        username: username,
        password,
      });

      if (result.status) {
        window.location.reload();
      }
    },
  };
};

const authContext = createContext<auth>({
  user: null,
  isAuthenticated: false,
  isChecking: true,
  signup: (username: string, password: string) => {},
  login: (username: string, password: string) => {},
  logout: () => {},
});

export const ProvideAuth: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const auth = useProvideAuth();
  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
};

export const useAuth = () => {
  return useContext(authContext);
};
