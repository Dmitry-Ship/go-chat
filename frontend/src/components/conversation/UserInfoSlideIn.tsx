import React from "react";
import styles from "./UserInfoSlideIn.module.css";
import Avatar from "../common/Avatar";
import SlideIn from "../common/SlideIn";
import { useAPI } from "../../contexts/apiContext";

const UserInfoSlideIn: React.FC<{
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  isOwner: boolean;
  isOpen: boolean;
  toggleUserInfo: () => void;
}> = ({ user, toggleUserInfo, isOpen, isOwner }) => {
  const { makeCommand } = useAPI();

  const handleChatClick = async (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    e.preventDefault();

    const result = await makeCommand("/startDirectConversation", {
      to_user_id: user.id,
    });

    if (result.status) {
      toggleUserInfo();

      window.location.href = `/conversations/${result.data.conversation_id}`;
    }
  };

  const handleKickClick = async (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    e.preventDefault();

    const result = await makeCommand("/kick", {
      user_id: user.id,
    });

    if (result.status) {
      toggleUserInfo();
    }
  };

  return (
    <SlideIn onClose={toggleUserInfo} isOpen={isOpen}>
      <div className={styles.userInfo}>
        <Avatar src={user.avatar} size={100} />
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
};

export default UserInfoSlideIn;
