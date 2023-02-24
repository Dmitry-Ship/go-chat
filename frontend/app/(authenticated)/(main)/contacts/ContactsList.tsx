"use client";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { Loader } from "../../../../src/components/common/Loader";
import { ContactItem } from "./ContactItem";
import styles from "./ContactsList.module.css";
import { EmptyScreen } from "../../../../src/components/common/EmptyScreen";
import { useMutation, useQuery } from "react-query";
import {
  getContacts,
  startDirectConversation,
} from "../../../../src/api/fetch";

export function ContactsList() {
  const [page, setPage] = useState(1);

  const { data, status } = useQuery({
    queryKey: ["contacts", page],
    queryFn: () => getContacts(page),
    keepPreviousData: true,
  });

  const handleScroll = (e: React.UIEvent<HTMLElement>) => {
    if (
      e.currentTarget.scrollHeight - e.currentTarget.scrollTop ===
      e.currentTarget.clientHeight
    ) {
      setPage(page + 1);
    }
  };

  const router = useRouter();

  const startDirectConversationRequest = useMutation(startDirectConversation, {
    onSuccess: (data) => {
      router.push(`/conversations/${data.conversation_id}`);
    },
  });

  const handleClick =
    (id: string) =>
    async (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
      e.preventDefault();
      startDirectConversationRequest.mutate({
        to_user_id: id,
      });
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
          switch (status) {
            case "loading":
              return <Loader />;
            case "success":
              return (
                <>
                  {data?.length ? (
                    data?.map((user, i) => (
                      <ContactItem
                        key={i}
                        onClick={handleClick(user.id)}
                        contact={user}
                      />
                    ))
                  ) : (
                    <EmptyScreen text="No contacts yet 🤷🏼" />
                  )}
                </>
              );
            default:
              return null;
          }
        })()}
      </main>
    </>
  );
}
