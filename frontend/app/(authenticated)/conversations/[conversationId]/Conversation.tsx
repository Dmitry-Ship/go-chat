"use client";
import React, { useEffect } from "react";
import styles from "./Conversation.module.css";
import { ChatForm } from "./ChatForm";
import { ChatLog } from "./ChatLog";
import { EditConversationBtn } from "./EditConversationBtn";
import { useWebSocket } from "../../../../src/contexts/WSContext";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Avatar } from "../../../../src/components/common/Avatar";
import { Loader } from "../../../../src/components/common/Loader";
import { ParticipantsList } from "./ParticipantsList";
import { useQuery } from "react-query";
import { getConversation } from "../../../../src/api/fetch";

export function Conversation({ conversationId }: { conversationId: string }) {
  const router = useRouter();
  const { onNotification } = useWebSocket();

  const { data, status, refetch } = useQuery(
    `conversation${conversationId}`,
    getConversation(`?conversation_id=${conversationId}`)
  );

  useEffect(() => {
    onNotification("conversation_deleted", (event) => {
      if (event.data.conversation_id === conversationId) {
        router.push("/");
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [conversationId, router]);

  const setIsJoined = () => {
    refetch();
  };

  useEffect(() => {
    onNotification("conversation_updated", (event) => {
      if (status === "success" && event.data.id === data.id) {
        refetch();
      }
    });

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data]);

  return (
    <>
      {(() => {
        switch (status) {
          case "error":
            return <div>Error</div>;
          case "loading":
            return <Loader />;
          case "success": {
            return (
              <>
                <header className={`header header-for-scrollable`}>
                  <Link href="/main" className={styles.backButton}>
                    👈
                  </Link>

                  <div className={styles.conversationInfo}>
                    <div className={styles.conversationGroupInfo}>
                      <Avatar src={data?.avatar} />
                      <h3 className={styles.conversationName}>{data?.name}</h3>
                    </div>

                    {data?.type === "group" && (
                      <ParticipantsList
                        conversationId={conversationId}
                        participantsCount={data.participants_count}
                      />
                    )}
                  </div>
                  {data?.type === "group" ? (
                    <EditConversationBtn
                      conversationId={conversationId}
                      conversation={data}
                      onLeave={() => setIsJoined()}
                    />
                  ) : (
                    <div />
                  )}
                </header>

                <section className="wrap">
                  <ChatLog
                    conversation={data}
                    isEmpty={data?.participants_count < 2}
                  />

                  <ChatForm
                    conversationId={conversationId}
                    conversationType={data?.type}
                    loading={false}
                    joined={data?.joined}
                    onJoin={() => setIsJoined()}
                  />
                </section>
              </>
            );
          }
          default:
            return null;
        }
      })()}
    </>
  );
}
