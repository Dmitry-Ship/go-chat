import React, { useEffect, useState } from "react";
import styles from "./EditConversationBtn.module.css";
import SlideIn from "../common/SlideIn";
import { makeCommand, makeQuery } from "../../api/fetch";
import { Contact } from "../../types/coreTypes";
import ContactItem from "../contacts/ContactItem";
import { useRouter } from "next/router";

const InviteMenu: React.FC<{}> = () => {
  const [isInviteMenuOpen, setIsInviteMenuOpen] = useState(false);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const router = useRouter();
  const conversationId = router.query.conversationId as string;
  const handleToggleInviteMenu = async () => {
    setIsInviteMenuOpen(!isInviteMenuOpen);
  };

  useEffect(() => {
    const queryContacts = async () => {
      const response = await makeQuery(
        `/getPotentialInvitees?conversation_id=${conversationId}`
      );
      if (response.status) {
        setContacts(response.data);
      }
    };

    if (isInviteMenuOpen) {
      queryContacts();
    }
  }, [isInviteMenuOpen, conversationId]);

  const handleClick =
    (id: string) =>
    async (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
      e.preventDefault();
      const result = await makeCommand("/inviteUserToConversation", {
        user_id: id,
        conversation_id: conversationId,
      });

      if (result.status) {
        setIsInviteMenuOpen(false);
      }
    };

  return (
    <>
      <button
        onClick={handleToggleInviteMenu}
        className={`btn ${styles.menuItem}`}
      >
        🤙 Invite
      </button>
      <SlideIn onClose={handleToggleInviteMenu} isOpen={isInviteMenuOpen}>
        {contacts.length > 0 ? (
          contacts.map((contact, i) => (
            <ContactItem
              key={i}
              onClick={handleClick(contact.id)}
              contact={contact}
            />
          ))
        ) : (
          <h3>No contacts to invite</h3>
        )}
      </SlideIn>
    </>
  );
};

export default InviteMenu;
