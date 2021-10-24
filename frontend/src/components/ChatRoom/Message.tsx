import React from "react";
import styles from "./Message.module.css";
import { Message } from "../../types/coreTypes";
import { useAuth } from "../../authContext";
import Avatar from "../common/Avatar";

const MessageComponent: React.FC<{
  message: Message;
  isFistInAGroup: boolean;
  isLastInAGroup: boolean;
}> = ({ message, isFistInAGroup, isLastInAGroup }) => {
  const date = new Date(message.createdAt * 1000);
  const time = `${date.getHours()}:${date.getMinutes()}`;

  const { user } = useAuth();

  const isOutbound = message.user.id === user?.id;
  const isSystem = message.type === "system";

  return (
    <div className={styles.message}>
      {isSystem ? (
        <div className={styles.systemMessage}>{message.text}</div>
      ) : (
        <>
          <div className={styles.avatarColumn}>
            {isLastInAGroup && <Avatar src={message.user.avatar} />}
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
