import { useRouter } from "next/router";
import React from "react";
import { makeCommand } from "../../api/fetch";
import { useQuery } from "../../api/hooks";
import { Contact } from "../../types/coreTypes";
import Loader from "../common/Loader";
import ContactItem from "./ContactItem";
import styles from "./ContactsList.module.css";

function ContactsList() {
  const response = useQuery<Contact[]>(`/getContacts`);

  const router = useRouter();
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
      <main className={`${styles.list} scrollable-content`}>
        {(() => {
          switch (response.status) {
            case "fetching":
              return <Loader />;
            case "done":
              return response.data?.map((user, i) => (
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
