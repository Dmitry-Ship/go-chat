import React from "react";
import { Conversation } from "./Conversation";

function ChatConversationPage({
  params,
}: {
  params: { conversationId: string };
}) {
  return <Conversation conversationId={params.conversationId} />;
}

export default ChatConversationPage;
