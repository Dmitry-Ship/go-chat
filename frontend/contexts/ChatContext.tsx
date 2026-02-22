"use client";

import { createContext, ReactNode, useCallback, useContext } from "react";
import { useRouter, useSelectedLayoutSegment } from "next/navigation";

interface ChatContextType {
  activeConversationId: string | null;
  setActiveConversation: (id: string | null) => void;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

export const ChatProvider = ({ children }: { children: ReactNode }) => {
  const router = useRouter();
  const selectedSegment = useSelectedLayoutSegment();
  const activeConversationId = selectedSegment ?? null;

  const setActiveConversation = useCallback((id: string | null) => {
    if (id === activeConversationId) {
      return;
    }

    if (!id) {
      router.push("/chat");
      return;
    }

    router.push(`/chat/${encodeURIComponent(id)}`);
  }, [activeConversationId, router]);

  return (
    <ChatContext.Provider value={{ activeConversationId, setActiveConversation }}>
      {children}
    </ChatContext.Provider>
  );
};

export const useChat = () => {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error("useChat must be used within a ChatProvider");
  }
  return context;
};
