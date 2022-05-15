import React from "react";
import styles from "./ContactItem.module.css";
import Avatar from "../common/Avatar";
import { Contact } from "../../types/coreTypes";

type ContactItemProps = {
  contact: Contact;
  onClick: (
    e: React.MouseEvent<HTMLAnchorElement, MouseEvent>
  ) => Promise<void>;
};

const ContactItem: React.FC<ContactItemProps> = ({ contact, onClick }) => {
  return (
    <a onClick={onClick} className={`${styles.contact} rounded`}>
      <Avatar src={contact.avatar} size={65} />
      <h3 className={styles.contactName}>{contact.name}</h3>
    </a>
  );
};

export default ContactItem;
