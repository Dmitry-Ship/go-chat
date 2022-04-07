import React from "react";
import styles from "./ConversationItem.module.css";
import Link from "next/link";
import Avatar from "../common/Avatar";

type ConversationItemProps = {
  name: string;
  href: string;
};

const ConversationItem: React.FC<ConversationItemProps> = ({ href, name }) => {
  return (
    <Link href={href}>
      <a className={`${styles.conversation} rounded`}>
        <Avatar src={"H"} size={55} />
        <h3 className={styles.conversationName}>{name}</h3>
      </a>
    </Link>
  );
};

export default ConversationItem;
