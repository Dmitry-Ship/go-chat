import React from "react";
import styles from "./ConversationItem.module.css";
import Link from "next/link";
import Avatar from "../common/Avatar";
import { Conversation } from "../../types/coreTypes";

type ConversationItemProps = {
  conversation: Conversation;
};

const ConversationItem: React.FC<ConversationItemProps> = ({
  conversation,
}) => {
  return (
    <Link href={"conversations/" + conversation.id}>
      <a className={`${styles.conversation} rounded`}>
        <Avatar src={conversation.avatar} size={65} />
        <h3 className={styles.conversationName}>{conversation.name}</h3>
      </a>
    </Link>
  );
};

export default ConversationItem;
