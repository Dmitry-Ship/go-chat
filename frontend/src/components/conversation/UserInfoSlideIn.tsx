import React from "react";
import styles from "./UserInfoSlideIn.module.css";
import Avatar from "../common/Avatar";
import SlideIn from "../common/SlideIn";
import { makeCommand } from "../../api/fetch";

const UserInfoSlideIn: React.FC<{
  user: {
    id: string;
    avatar: string;
    name: string;
  };
  isOpen: boolean;
  toggleUserInfo: () => void;
}> = ({ user, toggleUserInfo, isOpen }) => {
  const handleClick = async (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    e.preventDefault();

    const result = await makeCommand("/createPrivateConversationIfNotExists", {
      to_user_id: user.id,
    });

    if (result.status) {
      toggleUserInfo();

      window.location.href = `/conversations/${result.data.conversation_id}`;
    }
  };

  return (
    <SlideIn onClose={toggleUserInfo} isOpen={isOpen}>
      <div className={styles.userInfo}>
        <Avatar src={user.avatar} size={100} />
        <h3>{user.name}</h3>
        <button className={`btn`} onClick={handleClick}>
          💬 Chat
        </button>
      </div>
    </SlideIn>
  );
};

export default UserInfoSlideIn;
