import React, { createContext, useContext, useEffect, useState } from "react";
import { makeCommand, makeQuery } from "./api/fetch";
import { User } from "./types/coreTypes";

const ACCESS_TOKEN_LIFETIME = 1000 * 60 * 10;
const ACCESS_TOKEN_REFETCH_TIME = ACCESS_TOKEN_LIFETIME / 2;

type auth = {
  user: null | User;
  isAuthenticated: boolean;
  signup: (username: string, password: string) => void;
  login: (username: string, password: string) => void;
  logout: () => void;
};

export const useProvideAuth = (): auth => {
  const isPreviouslyAuthenticated = Boolean(
    localStorage.getItem("isAuthenticated")
  );
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<User | null>(null);

  const refreshToken = async () => {
    const result = await makeCommand("/refreshToken");

    if (!result.status) {
      setIsAuthenticated(false);
      localStorage.setItem("isAuthenticated", "");
    } else {
      setIsAuthenticated(true);
      localStorage.setItem("isAuthenticated", "true");
    }
  };

  useEffect(() => {
    if (isPreviouslyAuthenticated) {
      refreshToken();

      const interval = setInterval(refreshToken, ACCESS_TOKEN_REFETCH_TIME);

      return () => clearInterval(interval);
    }
  }, [isPreviouslyAuthenticated]);

  useEffect(() => {
    if (isAuthenticated) {
      const fetchUser = async () => {
        const result = await makeQuery("/getUser");

        if (result.status) {
          setUser(result.data);
        }
      };

      fetchUser();
    }
  }, [isAuthenticated]);

  return {
    isAuthenticated,
    user,
    login: async (username: string, password: string) => {
      const result = await makeCommand("/login", {
        user_name: username,
        password,
      });

      if (result.status) {
        localStorage.setItem("isAuthenticated", "true");
        window.location.reload();
      }
    },
    logout: async () => {
      const result = await makeCommand("/logout");

      if (result.status) {
        localStorage.setItem("isAuthenticated", "");
        window.location.reload();
      }
    },
    signup: async (username: string, password: string) => {
      const result = await makeCommand("/signup", {
        user_name: username,
        password,
      });

      if (result.status) {
        localStorage.setItem("isAuthenticated", "true");
        window.location.reload();
      }
    },
  };
};

const authContext = createContext<auth>({
  user: null,
  isAuthenticated: false,
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
