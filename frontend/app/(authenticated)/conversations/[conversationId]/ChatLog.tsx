import React, { useEffect, useRef, useState } from "react";
import styles from "./ChatLog.module.css";
import { ConversationFull, MessageRaw } from "../../../../src/types/coreTypes";
import MessageComponent from "./MessageComponent";
import Loader from "../../../../src/components/common/Loader";
import { usePaginatedQuery } from "../../../../src/api/hooks";
import { parseMessage } from "../../../../src/messages";
import { useWebSocket } from "../../../../src/contexts/WSContext";
import InviteMenu from "./InviteMenu";

const ChatLog: React.FC<{
  conversation: ConversationFull;
  isEmpty: boolean;
}> = ({ conversation, isEmpty }) => {
  const { onNotification } = useWebSocket();
  const [lastScrollHeight, setLastScrollHeight] = useState<number>(0);

  const [messagesQuery, append, loadNext] = usePaginatedQuery<MessageRaw>(
    `/getConversationsMessages?conversation_id=${conversation?.id}`,
    true
  );

  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const currentScroll =
      (containerRef.current?.scrollHeight || 0) - lastScrollHeight;

    containerRef.current?.scrollTo(0, currentScroll);
  }, [messagesQuery.items.length]);

  const handleScroll = (e: React.UIEvent<HTMLElement>) => {
    if (e.currentTarget.scrollTop === 0) {
      setLastScrollHeight(e.currentTarget.scrollHeight);
      loadNext();
    }
  };

  useEffect(() => {
    onNotification("message", (event) => {
      if (event.data.conversation_id === conversation.id) {
        append([event.data]);
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <main
      className={`${styles.log} scrollable-content inner-shadow`}
      onScroll={handleScroll}
      ref={containerRef}
    >
      {messagesQuery.status === "fetching" ? (
        <Loader />
      ) : (
        <>
          {isEmpty && (
            <div className={styles.emptyLog}>
              <div>
                <h4>It feels lonely here</h4>
                <InviteMenu conversationId={conversation.id} />
              </div>
            </div>
          )}

          {messagesQuery.items.length > 0 ? (
            messagesQuery.items.map(parseMessage).map((item, i) => {
              const previous = messagesQuery.items[i - 1];
              const isFistInAGroup =
                !previous ||
                previous?.type !== "text" ||
                (item.type === "text" && item?.user?.id !== previous?.user.id);

              const next = messagesQuery.items[i + 1];

              const isLastInAGroup =
                !next ||
                next?.type !== "text" ||
                (item.type === "text" && item.user.id !== next?.user.id);

              return (
                <MessageComponent
                  key={i}
                  conversation={conversation}
                  message={item}
                  isFistInAGroup={isFistInAGroup}
                  isLastInAGroup={isLastInAGroup}
                />
              );
            })
          ) : (
            <div className={styles.emptyLog}>
              <p>👋 No messages yet</p>
            </div>
          )}
        </>
      )}
    </main>
  );
};

export default ChatLog;
