import React from "react";
import { useQuery } from "../../api/hooks";
import { User } from "../../types/coreTypes";
import Loader from "../common/Loader";
import ConversationItem from "../conversations/ConversationItem";
import styles from "./ContactsList.module.css";

function ContactsList() {
  const response = useQuery<User[]>(`/getContacts`);

  return (
    <main className={`${styles.list} scrollable-content`}>
      {(() => {
        switch (response.status) {
          case "fetching":
            return <Loader />;
          case "done":
            return response.data?.map((user, i) => (
              <ConversationItem
                key={i}
                conversation={{
                  id: user.id,
                  name: user.name,
                  avatar: user.avatar,
                }}
              />
            ));
          default:
            return null;
        }
      })()}
    </main>
  );
}

export default ContactsList;
