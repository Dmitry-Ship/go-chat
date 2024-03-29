import React, { useReducer, useState } from "react";
import { SlideIn } from "../../../../src/components/common/SlideIn";
import styles from "./ParticipantsList.module.css";
import { ContactItem } from "../../(main)/contacts/ContactItem";
import { Loader } from "../../../../src/components/common/Loader";
import { InviteMenu } from "./InviteMenu";
import { useMutation, useQuery } from "react-query";
import {
  getParticipants,
  startDirectConversation,
} from "../../../../src/api/fetch";

export function ParticipantsList({
  participantsCount,
  conversationId,
}: {
  participantsCount: number;
  conversationId: string;
}) {
  const [isParticipantsListOpen, toggleParticipantsList] = useReducer(
    (open) => !open,
    false
  );
  const { data, status, refetch } = useQuery(
    "participants",
    getParticipants(`?conversation_id=${conversationId}`),
    {
      refetchOnWindowFocus: false,
      enabled: false,
    }
  );

  const startDirectConversationRequest = useMutation(startDirectConversation, {
    onSuccess: (data) => {
      handleTogglesParticipantsListOpen();
      window.location.href = `/conversations/${data.conversation_id}`;
    },
  });

  function handleTogglesParticipantsListOpen() {
    refetch();
    toggleParticipantsList();
  }

  function handleClick(id: string) {
    return async function (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) {
      e.preventDefault();

      startDirectConversationRequest.mutate({
        to_user_id: id,
      });
    };
  }

  return (
    <>
      <div
        className={styles.conversationParticipantsCount}
        onClick={handleTogglesParticipantsListOpen}
      >
        {participantsCount} participants
      </div>
      <SlideIn
        onClose={handleTogglesParticipantsListOpen}
        isOpen={isParticipantsListOpen}
      >
        {(() => {
          switch (status) {
            case "loading":
              return <Loader />;
            case "success":
              return (
                <>
                  <InviteMenu conversationId={conversationId} />
                  {data.map((contact, i) => (
                    <ContactItem
                      key={i}
                      onClick={handleClick(contact.id)}
                      contact={contact}
                    />
                  ))}
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
