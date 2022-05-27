import React, { useState } from "react";
import SlideIn from "../common/SlideIn";
import styles from "./ParticipantsList.module.css";
import { Contact } from "../../types/coreTypes";
import ContactItem from "../contacts/ContactItem";
import { useRouter } from "next/router";
import { useQueryOnDemand } from "../../api/hooks";
import Loader from "../common/Loader";
import { useAPI } from "../../contexts/apiContext";
import InviteMenu from "./InviteMenu";

const ParticipantsList: React.FC<{ participantsCount: number }> = ({
  participantsCount,
}) => {
  const [isParticipantsListOpen, setIsParticipantsListOpen] = useState(false);
  const router = useRouter();
  const conversationId = router.query.conversationId as string;
  const [response, load] = useQueryOnDemand<Contact[]>(
    `/getParticipants?conversation_id=${conversationId}`
  );

  const { makeCommand } = useAPI();

  const handleTogglesParticipantsListOpen = () => {
    load();
    setIsParticipantsListOpen(!isParticipantsListOpen);
  };

  const handleClick =
    (id: string) =>
    async (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
      e.preventDefault();

      const result = await makeCommand("/startDirectConversation", {
        to_user_id: id,
      });

      if (result.status) {
        handleTogglesParticipantsListOpen();

        window.location.href = `/conversations/${result.data.conversation_id}`;
      }
    };

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
          switch (response.status) {
            case "fetching":
              return <Loader />;
            case "done":
              return (
                <>
                  <InviteMenu />
                  {response.data.map((contact, i) => (
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
};

export default ParticipantsList;
