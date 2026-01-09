"use client";

import { useConversation } from "@/hooks/queries/useConversation";
import { useParticipants } from "@/hooks/queries/useParticipants";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { MoreVertical, Users, LogOut } from "lucide-react";
import { ParticipantsDialog } from "./ParticipantsDialog";
import { useState } from "react";

interface ChatHeaderProps {
  conversationId: string;
  onLeave?: () => void;
}

export const ChatHeader = ({ conversationId, onLeave }: ChatHeaderProps) => {
  const { data: conversation } = useConversation(conversationId);
  const { data: participants = [] } = useParticipants(conversationId);
  const [showParticipants, setShowParticipants] = useState(false);

  if (!conversation) {
    return null;
  }

  return (
    <>
      <div className="h-16 border-b flex items-center px-4 justify-between">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-semibold">
            {conversation.avatar}
          </div>
          <div>
            <h2 className="font-semibold">{conversation.name}</h2>
            <p className="text-sm text-gray-500">
              {participants.length} participant{participants.length !== 1 ? "s" : ""}
            </p>
          </div>
        </div>

        <DropdownMenu>
          <DropdownMenuTrigger>
            <Button size="icon" variant="ghost">
              <MoreVertical className="h-5 w-5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {conversation.type === "group" && (
              <>
                <DropdownMenuItem onClick={() => setShowParticipants(true)}>
                  <Users className="mr-2 h-4 w-4" />
                  Participants
                </DropdownMenuItem>
                <DropdownMenuSeparator />
              </>
            )}
            {conversation.joined && (
              <DropdownMenuItem onClick={onLeave} className="text-red-600">
                <LogOut className="mr-2 h-4 w-4" />
                Leave Conversation
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <ParticipantsDialog
        open={showParticipants}
        onOpenChange={setShowParticipants}
        conversationId={conversationId}
      />
    </>
  );
};
