import React, { useEffect } from "react";
import { Conversation } from "../../types/coreTypes";
import styles from "./Conversation.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { useQuery } from "../../api/hooks";
import EditConversationBtn from "./EditConversationBtn";
import { useWebSocket } from "../../contexts/WSContext";
import { useRouter } from "next/router";
import Link from "next/link";
import Avatar from "../common/Avatar";
import Loader from "../common/Loader";
import ParticipantsList from "./ParticipantsList";

const Conversation: React.FC = () => {
  const router = useRouter();
  const conversationId = router.query.conversationId as string;
  const { onNotification } = useWebSocket();

  const [conversationQuery, updateConversation] = useQuery<
    Conversation & {
      joined: boolean;
      participants_count: number;
    }
  >(`/getConversation?conversation_id=${conversationId}`);

  useEffect(() => {
    onNotification("conversation_deleted", (event) => {
      if (event.data.conversation_id === conversationId) {
        router.push("/");
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [conversationId, router]);

  const setIsJoined = (isJoined: boolean) => {
    updateConversation({
      joined: isJoined,
    });
  };

  useEffect(() => {
    onNotification("conversation_renamed", (event) => {
      if (
        conversationQuery.status === "done" &&
        event.data.conversation_id === conversationQuery.data.id
      ) {
        updateConversation({
          ...conversationQuery.data,
          name: event.data.new_name,
        });
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [conversationQuery]);

  if (conversationQuery.status === "fetching") {
    return <Loader />;
  }

  if (conversationQuery.status === "error") {
    return <div>Error</div>;
  }

  const conversation = conversationQuery.data;

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <Link href="/">
          <a className={styles.backButton}>👈</a>
        </Link>
        <div className={styles.conversationInfo}>
          <div className={styles.conversationGroupInfo}>
            <Avatar src={conversation.avatar} />
            <h3 className={styles.conversationName}>{conversation.name}</h3>
          </div>

          {conversation?.type === "group" && (
            <ParticipantsList
              participantsCount={conversationQuery.data.participants_count}
            />
          )}
        </div>
        {conversation?.type === "group" ? (
          <EditConversationBtn
            conversationId={conversationId}
            joined={conversationQuery.data.joined}
            onLeave={() => setIsJoined(false)}
          />
        ) : (
          <div />
        )}
      </header>

      <section className="wrap">
        <ChatLog conversationId={conversationId} />

        <ChatForm
          conversationId={conversationId}
          conversationType={conversation.type}
          loading={false}
          joined={conversationQuery.data.joined}
          onJoin={() => setIsJoined(true)}
        />
      </section>
    </>
  );
};

export default Conversation;
