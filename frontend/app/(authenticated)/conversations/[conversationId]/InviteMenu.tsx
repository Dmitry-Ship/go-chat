import React, { useState } from "react";
import { SlideIn } from "../../../../src/components/common/SlideIn";
import { ContactItem } from "../../(main)/contacts/ContactItem";
import { Loader } from "../../../../src/components/common/Loader";
import { useMutation, useQuery } from "react-query";
import {
  getPotentialInvitees,
  inviteUserToConversation,
} from "../../../../src/api/fetch";

export const InviteMenu: React.FC<{ conversationId: string }> = ({
  conversationId,
}) => {
  const [isInviteMenuOpen, setIsInviteMenuOpen] = useState(false);
  const { data, status, refetch } = useQuery(
    "invitees",
    getPotentialInvitees(`?conversation_id=${conversationId}`),
    {
      refetchOnWindowFocus: false,
      enabled: false,
    }
  );

  const inviteUserToConversationRequest = useMutation(
    inviteUserToConversation,
    {
      onSuccess: (data) => {
        setIsInviteMenuOpen(false);
      },
    }
  );

  const handleToggleInviteMenu = () => {
    refetch();
    setIsInviteMenuOpen(!isInviteMenuOpen);
  };

  const handleClick =
    (id: string) => (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
      e.preventDefault();
      inviteUserToConversationRequest.mutate({
        user_id: id,
        conversation_id: conversationId,
      });
    };

  return (
    <>
      <button onClick={handleToggleInviteMenu} className={`btn m-b`}>
        🤙 Invite
      </button>
      <SlideIn onClose={handleToggleInviteMenu} isOpen={isInviteMenuOpen}>
        {(() => {
          switch (status) {
            case "loading":
              return <Loader />;
            case "success":
              return (
                <>
                  {data.length > 0 ? (
                    data.map((contact, i) => (
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
