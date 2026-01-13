"use client";

import { MessageDTO, UserDTO } from "@/lib/types";

interface MessageItemProps {
  message: MessageDTO;
  users: Record<string, UserDTO>;
  isCurrentUser: boolean;
}

export const MessageItem = ({ message, users, isCurrentUser }: MessageItemProps) => {
  const user = users[message.user_id];
  const isSystemMessage = message.type === "system";

  if (isSystemMessage) {
    return (
      <div className="flex justify-center my-2">
        <span className="text-xs text-gray-500 italic px-2 py-1 bg-gray-50 rounded">
          {message.text}
        </span>
      </div>
    );
  }

  return (
    <div className={`flex mb-4 ${isCurrentUser ? "justify-end" : "justify-start"}`}>
      {!isCurrentUser && (
        <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white text-sm font-semibold mr-2 flex-shrink-0">
          {user?.avatar}
        </div>
      )}
      <div
        className={`max-w-[70%] px-4 py-2 rounded-lg ${
          isCurrentUser
            ? "bg-blue-500 text-white"
            : "bg-gray-100 text-gray-900"
        }`}
      >
        {!isCurrentUser && (
          <div className="text-xs font-semibold mb-1 opacity-70">
            {user?.name}
          </div>
        )}
        <div className="break-words">{message.text}</div>
        <div
          className={`text-xs mt-1 ${isCurrentUser ? "text-blue-100" : "text-gray-500"}`}
        >
          {new Date(message.created_at).toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
          })}
        </div>
      </div>
      {isCurrentUser && (
        <div className="w-8 h-8 rounded-full bg-green-500 flex items-center justify-center text-white text-sm font-semibold ml-2 flex-shrink-0">
          {user?.avatar}
        </div>
      )}
    </div>
  );
};
