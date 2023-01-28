import React, { useState } from "react";
import SlideIn from "../../../../src/components/common/SlideIn";
import { Contact } from "../../../../src/types/coreTypes";
import ContactItem from "../../(main)/contacts/ContactItem";
import { useQueryOnDemand } from "../../../../src/api/hooks";
import Loader from "../../../../src/components/common/Loader";
import { useAPI } from "../../../../src/contexts/apiContext";

const InviteMenu: React.FC<{ conversationId: string }> = ({
  conversationId,
}) => {
  const [isInviteMenuOpen, setIsInviteMenuOpen] = useState(false);
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
      <button onClick={handleToggleInviteMenu} className={`btn m-b`}>
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
