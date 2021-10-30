import React from "react";
import Chat from "../../components/chatRoom/Chat";
import LoggedInLayout from "../../components/common/LoggedInLayout";

function Index() {
  return (
    <LoggedInLayout>
      <Chat />
    </LoggedInLayout>
  );
}

export default Index;
