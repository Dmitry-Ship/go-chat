"use client";

import React, { useState } from "react";
import { NewConversationBtn } from "./NewConversationBtn";
import { Loader } from "../../../../src/components/common/Loader";
import { ConversationItem } from "./ConversationItem";
import { EmptyScreen } from "../../../../src/components/common/EmptyScreen";
import { useQuery } from "react-query";
import { getConversations } from "../../../../src/api/fetch";

export function Conversations() {
  const [page, setPage] = useState(1);

  const { data, status } = useQuery({
    queryKey: ["conversation", page],
    queryFn: () => getConversations(page),
    keepPreviousData: true,
  });

  const handleScroll = (e: React.UIEvent<HTMLElement>) => {
    if (
      e.currentTarget.scrollHeight - e.currentTarget.scrollTop ===
      e.currentTarget.clientHeight
    ) {
      setPage(page + 1);
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
            switch (status) {
              case "loading":
                return <Loader />;
              case "success": {
                return data.length === 0 ? (
                  <EmptyScreen text="No one to talk to yet 🤷🏼">
                    <NewConversationBtn text={"+ New Group Chat"} />
                  </EmptyScreen>
                ) : (
                  <>
                    {data?.map((conversation, i) => (
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
