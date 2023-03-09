import React, { useReducer } from "react";
import { SlideIn } from "../../../../src/components/common/SlideIn";
import { ContactItem } from "../../(main)/contacts/ContactItem";
import { Loader } from "../../../../src/components/common/Loader";
import { useMutation, useQuery } from "react-query";
import {
  getPotentialInvitees,
  inviteUserToConversation,
} from "../../../../src/api/fetch";

export function InviteMenu({ conversationId }: { conversationId: string }) {
  const [isInviteMenuOpen, toggleInviteMenu] = useReducer(
    (open) => !open,
    false
  );
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
        toggleInviteMenu();
      },
    }
  );

  function handleToggleInviteMenu() {
    refetch();
    toggleInviteMenu();
  }

  function handleClick(id: string) {
    return function (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) {
      e.preventDefault();
      inviteUserToConversationRequest.mutate({
        user_id: id,
        conversation_id: conversationId,
      });
    };
  }

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
}
