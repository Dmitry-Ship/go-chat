"use client";

import { forwardRef, useEffect } from "react";
import { MessageItem } from "./MessageItem";
import { MessageDTO } from "@/lib/types";

interface MessageListProps {
  messages: MessageDTO[];
  currentUserId: string;
  messagesEndRef: React.RefObject<HTMLDivElement | null>;
}

export const MessageList = forwardRef<HTMLDivElement, MessageListProps>(
  ({ messages, currentUserId, messagesEndRef }, ref) => {
    useEffect(() => {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages, messagesEndRef]);

    return (
      <div ref={ref} className="flex-1 overflow-y-auto p-4">
        <div className="space-y-4">
          {messages.length === 0 ? (
            <div className="text-center text-gray-500 mt-10">
              No messages yet. Start the conversation!
            </div>
          ) : (
            messages.map((message) => (
              <MessageItem
                key={message.id}
                message={message}
                isCurrentUser={!message.is_inbound}
              />
            ))
          )}
          <div ref={messagesEndRef} />
        </div>
      </div>
    );
  }
);

MessageList.displayName = "MessageList";
