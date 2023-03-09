"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { AuthenticationService, IAuthenticationService } from "../auth";
import { User } from "../types/coreTypes";

type auth = {
  user: null | User;
  isChecking: boolean;
  signup: (username: string, password: string) => void;
  login: (username: string, password: string) => void;
  logout: () => void;
};

export const useProvideAuth = (
  authenticationService: IAuthenticationService
): auth => {
  const [isChecking, setIsChecking] = useState<boolean>(true);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    authenticationService.onStateChanged((fetchedUser, error) => {
      setUser(fetchedUser);

      // setError(error);
      setIsChecking(false);
    });
    authenticationService.rotateTokens();

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return {
    isChecking,
    user,
    login: authenticationService.login,
    logout: authenticationService.logout,
    signup: authenticationService.signup,
  };
};

const authContext = createContext<auth>({
  user: null,
  isChecking: true,
  signup: (username: string, password: string) => {},
  login: (username: string, password: string) => {},
  logout: () => {},
});

const authenticationService = new AuthenticationService();

export const ProvideAuth = ({ children }: { children: React.ReactNode }) => {
  const auth = useProvideAuth(authenticationService);

  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
};

export const useAuth = () => {
  return useContext(authContext);
};
