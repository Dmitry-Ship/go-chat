import React from "react";
import styles from "./MessageComponent.module.css";
import { Message } from "../../types/coreTypes";
import Avatar from "../common/Avatar";

const MessageComponent: React.FC<{
  message: Message;
  isFistInAGroup: boolean;
  isLastInAGroup: boolean;
}> = ({ message, isFistInAGroup, isLastInAGroup }) => {
  const date = new Date(message.createdAt);
  const time = `${date.getHours()}:${date.getMinutes()}`;

  return (
    <div className={styles.message}>
      {(() => {
        switch (message.type) {
          case "text":
            return (
              <>
                <div className={styles.avatarColumn}>
                  {isLastInAGroup && <Avatar src={message.user.avatar} />}
                </div>

                <div
                  className={`${
                    message.isInbound
                      ? styles.inboundMessage
                      : styles.outboundMessage
                  } ${styles.messageBubble}`}
                >
                  {isFistInAGroup && (
                    <div className={styles.userName}>{message.user.name}</div>
                  )}
                  {message.text}
                  <div className={styles.time}>{time}</div>
                </div>
              </>
            );
          default:
            return <div className={styles.systemMessage}>{message.text}</div>;
        }
      })()}
    </div>
  );
};

export default MessageComponent;
