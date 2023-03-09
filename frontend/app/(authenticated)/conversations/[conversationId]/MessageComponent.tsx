import React, { useReducer } from "react";
import styles from "./MessageComponent.module.css";
import { ConversationFull, Message } from "../../../../src/types/coreTypes";
import { Avatar } from "../../../../src/components/common/Avatar";
import { UserInfoSlideIn } from "./UserInfoSlideIn";

export const MessageComponent = ({
  message,
  isFistInAGroup,
  isLastInAGroup,
  conversation,
}: {
  message: Message;
  conversation: ConversationFull;
  isFistInAGroup: boolean;
  isLastInAGroup: boolean;
}) => {
  const [isOpen, toggleOpen] = useReducer((open) => !open, false);

  const date = new Date(message.createdAt);
  const time = `${date.getHours()}:${date.getMinutes()}`;

  return (
    <div className={styles.message}>
      {(() => {
        switch (message.type) {
          case "text":
            return (
              <div
                className={`${
                  message.isInbound
                    ? styles.inboundMessage
                    : styles.outboundMessage
                } ${styles.textMessage}`}
              >
                <UserInfoSlideIn
                  toggleUserInfo={toggleOpen}
                  isOwner={conversation.is_owner}
                  user={message.user}
                  isOpen={isOpen}
                />
                {message.isInbound && (
                  <div className={styles.avatarColumn} onClick={toggleOpen}>
                    {isLastInAGroup && <Avatar src={message.user.avatar} />}
                  </div>
                )}

                <div className={`${styles.messageBubble} shadow`}>
                  {message.isInbound && isFistInAGroup && (
                    <div className={styles.userName} onClick={toggleOpen}>
                      {message.user.name}
                    </div>
                  )}

                  {message.text}
                  <div className={styles.time}>{time}</div>
                </div>
              </div>
            );
          default:
            return (
              <div className={styles.systemMessage}>
                <span>{message.text}</span>
              </div>
            );
        }
      })()}
    </div>
  );
};
