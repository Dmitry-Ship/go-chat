import React, { FormEvent, useState } from "react";
import styles from "./ChatForm.module.css";
import Loader from "../common/Loader";
import { useWebSocket } from "../../contexts/WSContext";
import { makeCommand } from "../../api/fetch";

const ChatForm: React.FC<{
  loading: boolean;
  joined: boolean;
  conversationId: string;
  onJoin: () => void;
}> = ({ loading, joined, onJoin, conversationId }) => {
  const [message, setMessage] = useState<string>("");

  const { sendNotification } = useWebSocket();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    sendNotification("message", {
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
          {joined ? (
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
