"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { makeCommand, makePaginatedQuery, makeQuery } from "../api/fetch";

type api = {
  makeCommand: (
    url: string,
    body?: Record<string, any>
  ) => Promise<{ status: boolean; data: any }>;
  makeQuery: (url: string) => Promise<{ status: boolean; data: any }>;
  makePaginatedQuery: (
    url: string,
    page: number,
    perPage?: number
  ) => Promise<{ status: boolean; data: any; nextPage: number }>;
  setError: (error: string | null) => void;
  error: string | null;
};

const apiContext = createContext<api | null>(null);

export const ProvideAPI: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (error) {
      setTimeout(() => {
        setError(null);
      }, 5000);

      return;
    }
  }, [error]);

  const makeCommandMethod = async (
    url: string,
    body?: Record<string, any> | undefined
  ) => {
    const response = await makeCommand(url, body);
    if (response.status === false) {
      setError(response.error);
    }

    return response;
  };

  const makeQueryMethod = async (url: string) => {
    const response = await makeQuery(url);
    if (response.status === false) {
      setError(response.error);
    }

    return response;
  };

  const makePaginatedQueryMethod = async (
    url: string,
    page: number,
    perPage?: number
  ) => {
    const response = await makePaginatedQuery(url, page, perPage);

    return response;
  };

  const api = {
    makeCommand: makeCommandMethod,
    makeQuery: makeQueryMethod,
    makePaginatedQuery: makePaginatedQueryMethod,
    setError: setError,
    error: error,
  };

  return <apiContext.Provider value={api}>{children}</apiContext.Provider>;
};

export const useAPI = () => {
  return useContext(apiContext) as api;
};
