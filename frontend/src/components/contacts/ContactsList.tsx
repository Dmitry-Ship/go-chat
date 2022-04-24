import React from "react";
import { useQuery } from "../../api/hooks";
import { Contact } from "../../types/coreTypes";
import Loader from "../common/Loader";
import ContactItem from "./ContactItem";
import styles from "./ContactsList.module.css";

function ContactsList() {
  const response = useQuery<Contact[]>(`/getContacts`);

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
                <ContactItem key={i} contact={user} />
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
