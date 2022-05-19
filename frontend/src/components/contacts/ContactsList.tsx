import { useRouter } from "next/router";
import React from "react";
import { usePaginatedQuery } from "../../api/hooks";
import { useAPI } from "../../contexts/apiContext";
import { Contact } from "../../types/coreTypes";
import Loader from "../common/Loader";
import ContactItem from "./ContactItem";
import styles from "./ContactsList.module.css";

function ContactsList() {
  const [contactsQuery, , loadNext] =
    usePaginatedQuery<Contact>("/getContacts");

  const handleScroll = (e: React.UIEvent<HTMLElement>) => {
    if (
      e.currentTarget.scrollHeight - e.currentTarget.scrollTop ===
      e.currentTarget.clientHeight
    ) {
      loadNext();
    }
  };

  const router = useRouter();

  const { makeCommand } = useAPI();
  const handleClick =
    (id: string) =>
    async (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
      e.preventDefault();

      const result = await makeCommand(
        "/createPrivateConversationIfNotExists",
        {
          to_user_id: id,
        }
      );

      if (result.status) {
        router.push(`/conversations/${result.data.conversation_id}`);
      }
    };

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <h2>Contacts</h2>
      </header>
      <main
        className={`${styles.list} scrollable-content`}
        onScroll={handleScroll}
      >
        {(() => {
          switch (contactsQuery.status) {
            case "fetching":
              return <Loader />;
            case "done":
              return contactsQuery.items?.map((user, i) => (
                <ContactItem
                  key={i}
                  onClick={handleClick(user.id)}
                  contact={user}
                />
              ));
            default:
              return null;
          }
        })()}
      </main>
    </>
  );
}

export default ContactsList;
