import React, { useEffect, useRef, useState } from "react";
import styles from "./ChatLog.module.css";
import { ConversationFull, MessageRaw } from "../../../../src/types/coreTypes";
import { MessageComponent } from "./MessageComponent";
import { Loader } from "../../../../src/components/common/Loader";
import { parseMessage } from "../../../../src/messages";
import { useWebSocket } from "../../../../src/contexts/WSContext";
import { InviteMenu } from "./InviteMenu";
import { useQuery } from "react-query";
import { getConversationsMessages } from "../../../../src/api/fetch";

export const ChatLog = ({
  conversation,
  isEmpty,
}: {
  conversation: ConversationFull;
  isEmpty: boolean;
}) => {
  const { onNotification } = useWebSocket();
  const [lastScrollHeight, setLastScrollHeight] = useState<number>(0);

  const [page, setPage] = useState(1);
  const [localMessages, setLocalMessages] = useState<MessageRaw[]>([]);

  const { data, status } = useQuery({
    queryKey: [`${conversation?.id}/messages`, page],
    queryFn: async () => {
      const result = await getConversationsMessages(
        page,
        `?conversation_id=${conversation?.id}`
      );
      setLocalMessages(result);
      return result;
    },
    keepPreviousData: true,
  });

  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const currentScroll =
      (containerRef.current?.scrollHeight || 0) - lastScrollHeight;

    containerRef.current?.scrollTo(0, currentScroll);
  }, [data?.length, lastScrollHeight]);

  const handleScroll = (e: React.UIEvent<HTMLElement>) => {
    if (e.currentTarget.scrollTop === 0) {
      setLastScrollHeight(e.currentTarget.scrollHeight);
      setPage(page + 1);
    }
  };

  useEffect(() => {
    onNotification("message", (event) => {
      if (event.data.conversation_id === conversation.id) {
        setLocalMessages([event.data, ...localMessages]);
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
      {(() => {
        switch (status) {
          case "loading":
            return <Loader />;
          case "success":
            return (
              <>
                {isEmpty && (
                  <div className={styles.emptyLog}>
                    <div>
                      <h4>It feels lonely here</h4>
                      <InviteMenu conversationId={conversation.id} />
                    </div>
                  </div>
                )}

                {localMessages.length > 0 ? (
                  localMessages.map(parseMessage).map((item, i) => {
                    const previous = data[i - 1];
                    const isFistInAGroup =
                      !previous ||
                      previous?.type !== "text" ||
                      (item.type === "text" &&
                        item?.user?.id !== previous?.user.id);

                    const next = data[i + 1];

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
            );
          default:
            return null;
        }
      })()}
    </main>
  );
};
