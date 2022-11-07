import React, { FormEvent, useState } from "react";
import styles from "./ChatForm.module.css";
import Loader from "../../../../src/components/common/Loader";
import { useWebSocket } from "../../../../src/contexts/WSContext";
import { useAPI } from "../../../../src/contexts/apiContext";

const ChatForm: React.FC<{
  loading: boolean;
  joined: boolean;
  conversationType: "group" | "direct";
  conversationId: string;
  onJoin: () => void;
}> = ({ loading, joined, onJoin, conversationId, conversationType }) => {
  const [message, setMessage] = useState<string>("");

  const { sendNotification } = useWebSocket();
  const { makeCommand } = useAPI();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const notification =
      conversationType === "group" ? "group_message" : "direct_message";

    sendNotification(notification, {
      content: message,
      conversation_id: conversationId,
    });

    setMessage("");
  };

  const handleJoin = async () => {
    await makeCommand("/joinConversation", { conversation_id: conversationId });
    onJoin();
  };

  return (
    <div className={"controls-for-scrollable"}>
      {loading ? (
        <Loader />
      ) : (
        <>
          {joined || conversationType === "direct" ? (
            <form className={styles.form} onSubmit={handleSubmit}>
              <input
                type="text"
                className={"input"}
                size={64}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
              />
              <button
                disabled={!message}
                type="submit"
                className={styles.submitBtn}
              >
                👌
              </button>
            </form>
          ) : (
            <button onClick={handleJoin} className="btn">
              Join
            </button>
          )}
        </>
      )}
    </div>
  );
};

export default ChatForm;
