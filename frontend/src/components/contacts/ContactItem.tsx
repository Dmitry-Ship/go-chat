import React from "react";
import styles from "./ContactItem.module.css";
import Avatar from "../common/Avatar";
import { Contact } from "../../types/coreTypes";
import { useRouter } from "next/router";
import { makeCommand } from "../../api/fetch";

type ContactItemProps = {
  contact: Contact;
};

const ContactItem: React.FC<ContactItemProps> = ({ contact }) => {
  const router = useRouter();
  const handleClick = async (
    e: React.MouseEvent<HTMLAnchorElement, MouseEvent>
  ) => {
    e.preventDefault();

    const result = await makeCommand("/createPrivateConversationIfNotExists", {
      to_user_id: contact.id,
    });

    if (result.status) {
      router.push(`/conversations/${result.data.conversation_id}`);
    }
  };

  return (
    <a onClick={handleClick} className={`${styles.contact} rounded`}>
      <Avatar src={contact.avatar} size={55} />
      <h3 className={styles.contactName}>{contact.name}</h3>
    </a>
  );
};

export default ContactItem;
