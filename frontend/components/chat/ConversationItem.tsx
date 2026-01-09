"use client";

import { ConversationDTO } from "@/lib/types";

interface ConversationItemProps {
  conversation: ConversationDTO;
  onClick: () => void;
  active: boolean;
}

export const ConversationItem = ({
  conversation,
  onClick,
  active,
}: ConversationItemProps) => {
  return (
    <div
      onClick={onClick}
      className={`flex items-center gap-3 p-3 cursor-pointer rounded-lg transition-colors ${
        active
          ? "bg-blue-100 dark:bg-blue-900"
          : "hover:bg-gray-100 dark:hover:bg-gray-800"
      }`}
    >
      <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-semibold flex-shrink-0">
        {conversation.avatar}
      </div>
      <div className="flex-1 min-w-0">
        <div className="font-medium truncate">{conversation.name}</div>
        {conversation.last_message && (
          <div className="text-sm text-gray-500 truncate">
            {conversation.last_message.text}
          </div>
        )}
      </div>
      {conversation.last_message && (
        <div className="text-xs text-gray-400 flex-shrink-0">
          {new Date(conversation.last_message.created_at).toLocaleDateString(
            [],
            { month: "short", day: "numeric" }
          )}
        </div>
      )}
    </div>
  );
};
