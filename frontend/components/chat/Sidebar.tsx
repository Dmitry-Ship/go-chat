"use client";

import { useState } from "react";
import { useChat } from "@/contexts/ChatContext";
import { useAuth } from "@/contexts/AuthContext";
import { useConversations } from "@/hooks/queries/useConversations";
import { ConversationItem } from "./ConversationItem";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Plus } from "lucide-react";
import { CreateGroupDialog } from "./CreateGroupDialog";
import { StartDirectDialog } from "./StartDirectDialog";

interface SidebarProps {
  className?: string;
  onConversationSelect?: (id: string) => void;
}

export const Sidebar = ({ className = "", onConversationSelect }: SidebarProps) => {
  const { setActiveConversation, activeConversationId } = useChat();
  const { user } = useAuth();
  const { data: conversations, isLoading } = useConversations();
  const [showCreateGroup, setShowCreateGroup] = useState(false);
  const [showStartDirect, setShowStartDirect] = useState(false);

  const handleConversationClick = (id: string) => {
    setActiveConversation(id);
    onConversationSelect?.(id);
  };

  return (
    <div className={`w-full md:w-80 border-r flex flex-col ${className}`}>
      <div className="p-4 border-b">
        <div className="flex items-center justify-between">
          <h1 className="text-xl font-bold">Go-Chat</h1>
          <div className="flex gap-2">
            <DropdownMenu>
              <DropdownMenuTrigger>
                <Button size="icon" variant="ghost">
                  <Plus className="h-5 w-5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => setShowCreateGroup(true)}>
                  New Group
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => setShowStartDirect(true)}>
                  New Direct Message
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        {isLoading ? (
          <div className="p-4 text-center text-gray-500">Loading...</div>
        ) : (
          conversations?.map((conversation) => (
            <ConversationItem
              key={conversation.id}
              conversation={conversation}
              active={conversation.id === activeConversationId}
              onClick={() => handleConversationClick(conversation.id)}
            />
          ))
        )}
      </div>

      {user && (
        <div className="p-4 border-t flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white font-semibold">
              {user.avatar}
            </div>
            <span className="font-medium">{user.name}</span>
          </div>
        </div>
      )}

      <CreateGroupDialog
        open={showCreateGroup}
        onOpenChange={setShowCreateGroup}
      />
      <StartDirectDialog
        open={showStartDirect}
        onOpenChange={setShowStartDirect}
      />
    </div>
  );
};
