import React from "react";
import { Conversation } from "../../types/coreTypes";
import { useQuery } from "../../api/hooks";
import NewConversationBtn from "./NewConversationBtn";
import Loader from "../common/Loader";
import ConversationItem from "./ConversationItem";

function Conversations() {
  const response = useQuery<Conversation[]>("/getConversations");

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <h2>Chats</h2>
        <NewConversationBtn />
      </header>
      <section className="wrap">
        <div className={`scrollable-content`}>
          {(() => {
            switch (response.status) {
              case "fetching":
                return <Loader />;
              case "done": {
                return response.data?.length === 0 ? (
                  <NewConversationBtn />
                ) : (
                  <>
                    {response.data?.map((conversation, i) => (
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
