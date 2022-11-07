"use client";
import React, { createContext, useContext, useEffect, useState } from "react";
import { AuthenticationService, IAuthenticationService } from "../auth";
import { User } from "../types/coreTypes";
import { useAPI } from "./apiContext";

type auth = {
  user: null | User;
  isAuthenticated: boolean;
  isChecking: boolean;
  signup: (username: string, password: string) => void;
  login: (username: string, password: string) => void;
  logout: () => void;
};

export const useProvideAuth = (
  authenticationService: IAuthenticationService
): auth => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isChecking, setIsChecking] = useState<boolean>(true);
  const [user, setUser] = useState<User | null>(null);
  const { setError } = useAPI();

  authenticationService.onLogin(() => {
    setIsAuthenticated(true);
    setIsChecking(false);
  });

  authenticationService.onLogout(() => {
    setIsAuthenticated(false);
    setIsChecking(false);
    setUser(null);
  });

  authenticationService.onError(setError);

  useEffect(() => {
    const getUser = async () => {
      const user = await authenticationService.fetchUser();
      setUser(user);
    };

    if (isAuthenticated) {
      getUser();
    }

    const timeout = authenticationService.rotateTokens();
    return () => clearTimeout(timeout);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  return {
    isAuthenticated,
    isChecking,
    user,
    login: authenticationService.login,
    logout: authenticationService.logout,
    signup: authenticationService.signup,
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
  const authenticationService = new AuthenticationService();
  const auth = useProvideAuth(authenticationService);
  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
};

export const useAuth = () => {
  return useContext(authContext);
};
