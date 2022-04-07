import React, { useEffect, useState } from "react";
import { Conversation } from "../../types/coreTypes";
import styles from "./Conversation.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { useQuery } from "../../api/hooks";
import EditConversationBtn from "./EditConversationBtn";
import { useWS } from "../../contexts/WSContext";
import { useRouter } from "next/router";
import Link from "next/link";

const Conversation: React.FC = () => {
  const router = useRouter();
  const conversationId = router.query.conversationId as string;
  const [conversation, setConversation] = useState<Conversation>();
  const [isJoined, setIsJoined] = useState(false);
  const { onNotification } = useWS();

  useEffect(() => {
    onNotification("conversation_deleted", (event) => {
      if (event.data.conversation_id === conversationId) {
        router.push("/");
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [conversationId, router]);

  const conversationQuery = useQuery<{
    conversation: Conversation;
    joined: boolean;
  }>(`/getConversation?conversation_id=${conversationId}`);

  useEffect(() => {
    if (conversationQuery.status === "done" && conversationQuery.data) {
      setConversation(conversationQuery.data.conversation);
      setIsJoined(conversationQuery.data.joined);
    }
  }, [conversationQuery]);

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <Link href="/">
          <a className={styles.backButton}>⏪</a>
        </Link>
        <b>{conversation?.name}</b>

        <EditConversationBtn
          conversationId={conversationId}
          joined={isJoined}
          onLeave={() => setIsJoined(false)}
        />
      </header>

      <section className="wrap">
        <ChatLog conversationId={conversationId} />

        <ChatForm
          conversationId={conversationId}
          loading={conversationQuery.status === "fetching"}
          joined={isJoined}
          onJoin={() => setIsJoined(true)}
        />
      </section>
    </>
  );
};

export default Conversation;
