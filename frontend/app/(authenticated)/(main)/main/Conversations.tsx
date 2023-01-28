"use client";

import React from "react";
import { ConversationListItem } from "../../../../src/types/coreTypes";
import { usePaginatedQuery } from "../../../../src/api/hooks";
import NewConversationBtn from "./NewConversationBtn";
import Loader from "../../../../src/components/common/Loader";
import ConversationItem from "./ConversationItem";
import EmptyScreen from "../../../../src/components/common/EmptyScreen";

function Conversations() {
  const [conversationsQuery, , loadNext] =
    usePaginatedQuery<ConversationListItem>("/getConversations");

  const handleScroll = (e: React.UIEvent<HTMLElement>) => {
    if (
      e.currentTarget.scrollHeight - e.currentTarget.scrollTop ===
      e.currentTarget.clientHeight
    ) {
      loadNext();
    }
  };

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <h2>Chats</h2>
        <NewConversationBtn text="+ New" />
      </header>
      <section className="wrap">
        <div className={`scrollable-content`} onScroll={handleScroll}>
          {(() => {
            switch (conversationsQuery.status) {
              case "fetching":
                return <Loader />;
              case "done": {
                return conversationsQuery.items.length === 0 ? (
                  <EmptyScreen text="No one to talk to yet 🤷🏼">
                    <NewConversationBtn text={"+ New Group Chat"} />
                  </EmptyScreen>
                ) : (
                  <>
                    {conversationsQuery.items?.map((conversation, i) => (
                      <ConversationItem key={i} conversation={conversation} />
                    ))}
                  </>
                );
              }
              default:
                return null;
            }
          })()}
        </div>
      </section>
    </>
  );
}

export default Conversations;
