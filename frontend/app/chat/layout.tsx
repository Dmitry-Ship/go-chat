"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { useAuth } from "@/contexts/AuthContext";
import { ChatProvider } from "@/contexts/ChatContext";
import { wsManager } from "@/lib/websocket";

export default function ChatLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { loading, authenticated } = useAuth();
  const router = useRouter();
  const queryClient = useQueryClient();

  useEffect(() => {
    wsManager.setQueryClient(queryClient);
  }, [queryClient]);

  useEffect(() => {
    if (!loading && !authenticated) {
      router.push("/login");
      return;
    }

    if (authenticated && !wsManager.isConnected()) {
      wsManager.connect();
    }

    return () => {
      wsManager.disconnect();
    };
  }, [loading, authenticated, router]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!authenticated) {
    return null;
  }

  return <ChatProvider>{children}</ChatProvider>;
}
