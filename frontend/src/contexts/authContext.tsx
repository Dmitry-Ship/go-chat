import React, { createContext, useContext, useEffect, useState } from "react";
import { login, logout, signup, rotateTokens, fetchUser } from "../auth";
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
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const getUser = async () => {
      const user = await fetchUser();
      setUser(user);
    };

    if (isAuthenticated) {
      getUser();
    }

    const authCallback = (status: boolean) => {
      setIsAuthenticated(status);
      setIsChecking(false);
    };

    const timeout = rotateTokens(authCallback);
    return () => clearTimeout(timeout);
  }, [isAuthenticated]);

  return {
    isAuthenticated,
    isChecking,
    user,
    login: (username: string, password: string) => {
      login(username, password, () => {
        setIsAuthenticated(true);
      });
    },
    logout: () => {
      logout(() => {
        setIsAuthenticated(false);
        setUser(null);
      });
    },
    signup: (username: string, password: string) => {
      signup(username, password, () => {
        setIsAuthenticated(true);
      });
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
