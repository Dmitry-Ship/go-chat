"use client";

import { useState } from "react";
import { useChat } from "@/contexts/ChatContext";
import { useContacts } from "@/hooks/queries/useContacts";
import { useStartDirectConversation } from "@/hooks/mutations/conversationMutations";
import {
  AlertDialog as Dialog,
  AlertDialogContent as DialogContent,
  AlertDialogDescription as DialogDescription,
  AlertDialogFooter as DialogFooter,
  AlertDialogHeader as DialogHeader,
  AlertDialogTitle as DialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { UserDTO } from "@/lib/types";

interface StartDirectDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export const StartDirectDialog = ({
  open,
  onOpenChange,
}: StartDirectDialogProps) => {
  const { setActiveConversation } = useChat();
  const { data: contacts = [] } = useContacts();
  const startDirectConversation = useStartDirectConversation();
  const [search, setSearch] = useState("");
  const [selectedUser, setSelectedUser] = useState<UserDTO | null>(null);

  const filteredContacts = contacts.filter((contact) =>
    contact.name.toLowerCase().includes(search.toLowerCase())
  );

  const handleStartConversation = async () => {
    if (!selectedUser) return;

    try {
      const result = await startDirectConversation.mutateAsync(selectedUser.id);
      setActiveConversation(result.conversation_id);
      setSelectedUser(null);
      onOpenChange(false);
    } catch (error) {
      console.error("Failed to start conversation:", error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Start Direct Message</DialogTitle>
          <DialogDescription>
            Select a user to start a direct conversation
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <Input
            placeholder="Search users..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />

          <div className="max-h-60 overflow-y-auto space-y-2">
            {filteredContacts.map((contact) => (
              <div
                key={contact.id}
                onClick={() => setSelectedUser(contact)}
                className={`flex items-center gap-3 p-2 rounded-lg cursor-pointer transition-colors ${
                  selectedUser?.id === contact.id
                    ? "bg-blue-100 dark:bg-blue-900"
                    : "hover:bg-gray-100 dark:hover:bg-gray-800"
                }`}
              >
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center text-white font-semibold">
                  {contact.avatar}
                </div>
                <div>
                  <div className="font-medium">{contact.name}</div>
                </div>
              </div>
            ))}
            {filteredContacts.length === 0 && (
              <div className="text-center text-gray-500 py-4">
                No contacts found
              </div>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={startDirectConversation.isPending}
          >
            Cancel
          </Button>
          <Button
            onClick={handleStartConversation}
            disabled={startDirectConversation.isPending || !selectedUser}
          >
            {startDirectConversation.isPending ? "Starting..." : "Start Chat"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
