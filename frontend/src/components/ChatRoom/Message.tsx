import React from "react";
import styles from "./Message.module.css";
import { Message } from "../../types/coreTypes";
import { useAuth } from "../../authContext";

const MessageComponent: React.FC<{
  message: Message;
  isFistInAGroup: boolean;
  isLastInAGroup: boolean;
}> = ({ message, isFistInAGroup, isLastInAGroup }) => {
  const date = new Date(message.createdAt * 1000);
  const time = `${date.getHours()}:${date.getMinutes()}`;

  const user = useAuth().user;

  const isOutbound = message.user.id === user?.id;
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
            className={`${
              isOutbound ? styles.outboundMessage : styles.inboundMessage
            } ${styles.messageBubble}`}
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
