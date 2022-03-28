import Link from "next/link";
import React from "react";
import { useQuery } from "../../api/hooks";
import { User } from "../../types/coreTypes";
import Avatar from "../common/Avatar";
import Loader from "../common/Loader";
import styles from "./ContactsList.module.css";

function ContactsList() {
  const response = useQuery<User[]>(`/getContacts`);

  return (
    <main className={`${styles.list} scrollable-content`}>
      {" "}
      {(() => {
        switch (response.status) {
          case "fetching":
            return <Loader />;
          case "done":
            return response.data?.map((user, i) => (
              <Link key={i} href={"rooms/" + user.id}>
                <a className={`${styles.contact} rounded`}>
                  <Avatar size={50} src={user.avatar} />
                  <h3>{user.name}</h3>
                </a>
              </Link>
            ));
          default:
            return null;
        }
      })()}
    </main>
  );
}

export default ContactsList;
