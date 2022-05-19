import React, { useState } from "react";
import styles from "./EditConversationBtn.module.css";
import SlideIn from "../common/SlideIn";
import { Contact } from "../../types/coreTypes";
import ContactItem from "../contacts/ContactItem";
import { useRouter } from "next/router";
import { useQueryOnDemand } from "../../api/hooks";
import Loader from "../common/Loader";
import { useAPI } from "../../contexts/apiContext";

const InviteMenu: React.FC<{}> = () => {
  const [isInviteMenuOpen, setIsInviteMenuOpen] = useState(false);
  const router = useRouter();
  const conversationId = router.query.conversationId as string;
  const [response, load] = useQueryOnDemand<Contact[]>(
    `/getPotentialInvitees?conversation_id=${conversationId}`
  );

  const { makeCommand } = useAPI();

  const handleToggleInviteMenu = () => {
    load();
    setIsInviteMenuOpen(!isInviteMenuOpen);
  };

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
        {(() => {
          switch (response.status) {
            case "fetching":
              return <Loader />;
            case "done":
              return (
                <>
                  {response.data.length > 0 ? (
                    response.data.map((contact, i) => (
                      <ContactItem
                        key={i}
                        onClick={handleClick(contact.id)}
                        contact={contact}
                      />
                    ))
                  ) : (
                    <h3>No contacts to invite</h3>
                  )}
                </>
              );
            default:
              return null;
          }
        })()}
      </SlideIn>
    </>
  );
};

export default InviteMenu;
