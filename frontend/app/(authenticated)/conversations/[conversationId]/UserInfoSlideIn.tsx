import React from "react";
import styles from "./UserInfoSlideIn.module.css";
import { Avatar } from "../../../../src/components/common/Avatar";
import { SlideIn } from "../../../../src/components/common/SlideIn";
import { useMutation } from "react-query";
import { kick, startDirectConversation } from "../../../../src/api/fetch";

export function UserInfoSlideIn({
  user,
  toggleUserInfo,
  isOpen,
  isOwner,
}: {
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  isOwner: boolean;
  isOpen: boolean;
  toggleUserInfo: () => void;
}) {
  const startDirectConversationRequest = useMutation(startDirectConversation, {
    onSuccess: (data) => {
      toggleUserInfo();

      window.location.href = `/conversations/${data.conversation_id}`;
    },
  });

  const kickRequest = useMutation(kick, {
    onSuccess: (data) => {
      toggleUserInfo();
    },
  });

  async function handleChatClick(
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) {
    e.preventDefault();

    startDirectConversationRequest.mutate({
      to_user_id: user.id,
    });
  }

  async function handleKickClick(
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) {
    e.preventDefault();

    kickRequest.mutate({
      user_id: user.id,
    });
  }

  return (
    <SlideIn onClose={toggleUserInfo} isOpen={isOpen}>
      <div className={styles.userInfo}>
        <Avatar src={user.avatar} size={"lg"} />
        <h3>{user.name}</h3>
        <button className={`btn m-b`} onClick={handleChatClick}>
          💬 Chat
        </button>

        {isOwner && (
          <button className={`btn`} onClick={handleKickClick}>
            🫵 Kick
          </button>
        )}
      </div>
    </SlideIn>
  );
}
