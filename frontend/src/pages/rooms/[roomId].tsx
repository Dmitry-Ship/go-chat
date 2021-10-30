import React from "react";
import ChatRoom from "../../components/chr/ChatRoom";
import LoggedInLayout from "../../components/common/LoggedInLayout";

function ChatRoomPage() {
  return (
    <LoggedInLayout>
      <ChatRoom />
    </LoggedInLayout>
  );
}

export default ChatRoomPage;
