import React from "react";
import styles from "./ConversationItem.module.css";
import Link from "next/link";
import Avatar from "../../../../src/components/common/Avatar";
import { ConversationListItem } from "../../../../src/types/coreTypes";

type ConversationItemProps = {
  conversation: ConversationListItem;
};

const ConversationItem: React.FC<ConversationItemProps> = ({
  conversation,
}) => {
  return (
    <Link
      href={"conversations/" + conversation.id}
      className={`${styles.wrap} rounded shadow`}
    >
      <>
        <Avatar src={conversation.avatar} />
        <div className={styles.conversationInfo}>
          <h3 className={styles.conversationName}>{conversation.name}</h3>

          {conversation.last_message && (
            <div className={styles.lastMessage}>
              {conversation.last_message.type === "text" &&
                conversation.type === "group" && (
                  <div>
                    <strong>{conversation.last_message.user.name}: </strong>
                  </div>
                )}
              <span className={styles.lastMessageText}>
                {conversation.last_message.text}
              </span>
            </div>
          )}
        </div>
      </>
    </Link>
  );
};

export default ConversationItem;
