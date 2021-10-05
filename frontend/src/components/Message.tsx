import React, { useContext } from "react";
import styles from "./Message.module.css";
import { Message } from "../types/coreTypes";
import { UserContext } from "../userContext";

const MessageComponent: React.FC<{
  message: Message;
  isFistInAGroup: boolean;
  isLastInAGroup: boolean;
}> = ({ message, isFistInAGroup, isLastInAGroup }) => {
  const date = new Date(message.created_at * 1000);
  const time = `${date.getHours()}:${date.getMinutes()}`;

  const user = useContext(UserContext);

  const isInbound = message.user.id === user.id;
  const isSystem = message.type === "system";

  return (
    <div className={styles.message}>
      {isSystem ? (
        <div className={styles.systemMessage}>{message.text}</div>
      ) : (
        <>
          <div className={styles.avatarColumn}>
            {isLastInAGroup && (
              <div className={styles.avatar}>{message.user.avatar}</div>
            )}
          </div>

          <div
            className={`${isInbound ? "" : styles.outboundMessage} ${
              styles.messageBubble
            }`}
          >
            {isFistInAGroup && (
              <div className={styles.userName}>{message.user.name}</div>
            )}
            {message.text}
            <div className={styles.time}>{time}</div>
          </div>
        </>
      )}
    </div>
  );
};

export default MessageComponent;
