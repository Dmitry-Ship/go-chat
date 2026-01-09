"use client";

import { useRef, useCallback } from "react";
import { useChat } from "@/contexts/ChatContext";
import { useConversation } from "@/hooks/queries/useConversation";
import { useMessages } from "@/hooks/queries/useMessages";
import { ChatHeader } from "./ChatHeader";
import { MessageList } from "./MessageList";
import { MessageInput } from "./MessageInput";

import { wsManager } from "@/lib/websocket";
import { useAuth } from "@/contexts/AuthContext";

interface ChatAreaProps {
  className?: string;
}

export const ChatArea = ({ className = "" }: ChatAreaProps) => {
  const { activeConversationId, setActiveConversation } = useChat();
  const { user } = useAuth();
  const { data: conversation } = useConversation(activeConversationId);
  const { data: messages } = useMessages(activeConversationId);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const isDirect = conversation?.type === "direct";

  const handleSendMessage = useCallback((content: string) => {
    if (!activeConversationId) return;

    const type = isDirect ? "direct_message" : "group_message";
    wsManager.sendMessage(type, { content, conversation_id: activeConversationId });
  }, [activeConversationId, isDirect]);

  if (!activeConversationId) {
    return (
      <div className={`flex items-center justify-center ${className}`}>
        <div className="text-center">
          <div className="w-24 h-24 mx-auto mb-4 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
            <svg
              className="w-12 h-12 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              />
            </svg>
          </div>
          <h2 className="text-xl font-semibold mb-2">Welcome to Go-Chat</h2>
          <p className="text-gray-500">
            Select a conversation to start chatting
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={`flex flex-col ${className}`}>
      <ChatHeader
        conversationId={activeConversationId}
        onLeave={() => setActiveConversation(null)}
      />

      <MessageList
        messages={messages || []}
        currentUserId={user?.id || ""}
        messagesEndRef={messagesEndRef}
      />

      <MessageInput
        onSend={handleSendMessage}
        disabled={false}
      />
    </div>
  );
};
