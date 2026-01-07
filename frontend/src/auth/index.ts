// import { getUser, login, logout, refreshToken, signup } from "../api/fetch";
import { useState } from "react";
import { useMutation, useQuery } from "react-query";
import * as api from "../api/fetch";
// import { User } from "../types/coreTypes";

let timeoutId: NodeJS.Timeout | null = null;

export function rotateTokens(): void {
  const refresh = async () => {
    const result = await api.refreshToken();
    if (result.access_token_expiration) {
      const accessTokenRefreshInterval =
        result.access_token_expiration / 1000000 / 2;
      timeoutId = setTimeout(refresh, accessTokenRefreshInterval);
    } else {
      clearTimeout(timeoutId!);
    }
  };

  if (!timeoutId) {
    timeoutId = setTimeout(refresh, 0);
  }
}

// async function fetchUser(): Promise<void> {
//   const getUserResult = await api.getUser()();
//   if (getUserResult) {
//     onStateChangedCallback(getUserResult, "");
//   }
// }

export function useAuth2() {
  const [isLoggenIn, setIsLoggedIn] = useState(false);
  const userQuery = useQuery(["user", isLoggenIn], api.getUser(), {
    enabled: isLoggenIn,
  });

  const login = useMutation(api.login, {
    onSuccess: () => {
      setIsLoggedIn(true);
    },
  });
  const signup = useMutation(api.signup, {
    onSuccess: () => {
      setIsLoggedIn(true);
    },
  });

  const logout = useMutation(api.logout, {
    onSuccess: () => {
      setIsLoggedIn(false);
    },
  });

  return {
    login: login.mutate,
    signup: signup.mutate,
    logout: logout.mutate,
    user: userQuery.data,
  };
}
