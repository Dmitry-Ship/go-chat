import React, { createContext, useContext, useEffect, useState } from "react";
import { makeRequest } from "./api/fetch";
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
    const result = await makeRequest("/refreshToken", {
      method: "POST",
    });

    if (!result.status) {
      setIsAuthenticated(false);
    } else {
      setIsAuthenticated(true);
    }
    setIsChecking(false);
  };

  useEffect(() => {
    refreshToken();

    if (isAuthenticated) {
      const fetchUser = async () => {
        const result = await makeRequest("/getUser");

        if (result.status) {
          setUser(result.data);
        }
      };

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
      const result = await makeRequest("/login", {
        method: "POST",
        body: {
          user_name: username,
          password,
        },
      });

      if (result.status) {
        window.location.replace("/");
      }
    },
    logout: async () => {
      const result = await makeRequest("/logout", {
        method: "POST",
      });

      if (result.status) {
        window.location.replace("/");
      }
    },
    signup: async (username: string, password: string) => {
      const result = await makeRequest("/signup", {
        method: "POST",
        body: {
          user_name: username,
          password,
        },
      });
      if (result.status) {
        window.location.replace("/");
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
