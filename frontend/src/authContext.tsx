import React, { createContext, useContext, useEffect, useState } from "react";
import { makeCommand, makeQuery } from "./api/fetch";
import { User } from "./types/coreTypes";

const ACCESS_TOKEN_LIFETIME = 1000 * 60 * 10;
const ACCESS_TOKEN_REFETCH_TIME = ACCESS_TOKEN_LIFETIME / 2;

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
  const [user, setUser] = useState<User | null>(null);

  const refreshToken = async () => {
    const result = await makeCommand("/refreshToken");

    if (!result.status) {
      setIsAuthenticated(false);
    } else {
      setIsAuthenticated(true);
    }
    setIsChecking(false);
  };

  const fetchUser = async () => {
    const result = await makeQuery("/getUser");

    if (result.status) {
      setUser(result.data);
    }
  };

  useEffect(() => {
    refreshToken();

    if (isAuthenticated) {
      fetchUser();
      const interval = setInterval(refreshToken, ACCESS_TOKEN_REFETCH_TIME);

      return () => clearInterval(interval);
    }
  }, [isAuthenticated]);

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
