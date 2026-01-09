"use client";

import { useEffect, useRef } from "react";
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getUser, login, signup, logout as apiLogout, refreshToken as apiRefreshToken } from '@/lib/api'

export function useAuth() {
  const queryClient = useQueryClient()
  const refreshTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const userQuery = useQuery({
    queryKey: ['auth', 'user'],
    queryFn: getUser,
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 15 * 60 * 1000,
  })

  const loginMutation = useMutation({
    mutationFn: ({ username, password }: { username: string; password: string }) =>
      login(username, password),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['auth'] })
      scheduleTokenRefresh(data.access_token_expiration);
    },
  })

  const signupMutation = useMutation({
    mutationFn: ({ username, password }: { username: string; password: string }) =>
      signup(username, password),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['auth'] })
      scheduleTokenRefresh(data.access_token_expiration);
    },
  })

  const logoutMutation = useMutation({
    mutationFn: apiLogout,
    onSuccess: () => {
      queryClient.clear()
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current);
        refreshTimeoutRef.current = null;
      }
    },
  })

  const refreshTokenMutation = useMutation({
    mutationFn: apiRefreshToken,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['auth'] })
      scheduleTokenRefresh(data.access_token_expiration);
    },
    onError: () => {
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current);
        refreshTimeoutRef.current = null;
      }
    },
  })

  const scheduleTokenRefresh = (expirationNanoseconds: number) => {
    if (refreshTimeoutRef.current) {
      clearTimeout(refreshTimeoutRef.current);
    }

    const expirationMs = expirationNanoseconds / 1e6;
    const refreshTime = expirationMs - 60 * 1000;

    if (refreshTime > 0) {
      refreshTimeoutRef.current = setTimeout(() => {
        refreshTokenMutation.mutate();
      }, refreshTime);
    }
  };

  useEffect(() => {
    return () => {
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current);
      }
    };
  }, []);

  return {
    user: userQuery.data || null,
    loading: userQuery.isLoading,
    authenticated: !!userQuery.data,
    login: loginMutation.mutateAsync,
    signup: signupMutation.mutateAsync,
    logout: logoutMutation.mutateAsync,
    refreshToken: refreshTokenMutation.mutateAsync,
    loginError: loginMutation.error,
    signupError: signupMutation.error,
    isLoggingIn: loginMutation.isPending,
    isSigningUp: signupMutation.isPending,
    isLoggingOut: logoutMutation.isPending,
  }
}
