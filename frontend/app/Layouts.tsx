"use client";
import React from "react";
import { ProvideAuth } from "../src/contexts/authContext";
import ErrorAlert from "./ErrorAlert";
import { QueryClient, QueryClientProvider } from "react-query";

const queryClient = new QueryClient();

export default function Layouts({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <ProvideAuth>
        <ErrorAlert />
        {children}
      </ProvideAuth>
    </QueryClientProvider>
  );
}
