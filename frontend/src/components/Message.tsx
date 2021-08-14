import React from "react";
import styles from "./Message.module.css";
import { Message } from "../types/coreTypes";

const classes = {
  system: styles.systemMessage,
  outbound: styles.outboundMessage,
  user: "",
};

const MessageComponent: React.FC<{ message: Message }> = ({ message }) => {
  const date = new Date(message.created_at);
  const time = `${date.getHours()}:${date.getMinutes()}`;

  return (
    <div className={styles.message}>
      {message.type === "system" ? (
        <div className={styles.systemMessage}>{message.text}</div>
      ) : (
        <>
          <div className={styles.avatar}>{message.sender.avatar}</div>
          <div className={`${classes[message.type]} ${styles.messageBubble}`}>
            <div className={styles.userName}>{message.sender.name}</div>
            {message.text}
            <div className={styles.time}>{time}</div>
          </div>
        </>
      )}
    </div>
  );
};

export default MessageComponent;
