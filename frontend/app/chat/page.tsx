"use client";

import { useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { Sidebar } from "@/components/chat/Sidebar";
import { ChatArea } from "@/components/chat/ChatArea";
import { Button } from "@/components/ui/button";
import { ArrowLeft, LogOut } from "lucide-react";
import { useChat } from "@/contexts/ChatContext";

export default function ChatPage() {
  const { user, logout } = useAuth();
  const { setActiveConversation } = useChat();
  const [showSidebar, setShowSidebar] = useState(true);

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="flex h-screen bg-white dark:bg-gray-900">
      <div className="flex-1 w-full md:hidden">
        <div className="h-14 border-b flex items-center justify-between px-4">
          <h1 className="font-bold text-lg">Go-Chat</h1>
          <div className="flex items-center gap-2">
            {user && (
              <div className="w-8 h-8 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white text-sm font-semibold">
                {user.avatar}
              </div>
            )}
            <Button size="icon" variant="ghost" onClick={handleLogout}>
              <LogOut className="h-5 w-5" />
            </Button>
          </div>
        </div>

        {showSidebar ? (
          <Sidebar
            className="h-[calc(100vh-3.5rem)]"
            onConversationSelect={(id) => {
              setActiveConversation(id);
              setShowSidebar(false);
            }}
          />
        ) : (
          <div className="h-[calc(100vh-3.5rem)] flex flex-col">
            <Button
              variant="ghost"
              onClick={() => setShowSidebar(true)}
              className="m-2 w-fit"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back
            </Button>
            <ChatArea className="flex-1" />
          </div>
        )}
      </div>

      <div className="hidden md:flex w-full h-screen">
        <Sidebar
          className="h-full"
          onConversationSelect={(id) => setActiveConversation(id)}
        />
        <ChatArea className="flex-1" />
      </div>
    </div>
  );
}
