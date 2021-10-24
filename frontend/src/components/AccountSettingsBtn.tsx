import React, { useState } from "react";
import styles from "./AccountSettingsBtn.module.css";
import SlideIn from "./common/SlideIn";
import { useAuth } from "../authContext";
import Avatar from "./common/Avatar";

const AccountSettingsBtn: React.FC = () => {
  const [isEditing, setIsEditing] = useState(false);
  const { user, logout } = useAuth();

  return (
    <>
      <button onClick={() => setIsEditing(true)} className={"navBtn"}>
        ⚙️
      </button>
      <SlideIn onClose={() => setIsEditing(false)} isOpen={isEditing}>
        <div>
          <div className={styles.accountInfo}>
            <Avatar src={user?.avatar || ""} />
            <h3>{user?.name}</h3>
          </div>
          <button onClick={logout} className={`btn`}>
            Logout
          </button>
        </div>
      </SlideIn>
    </>
  );
};

export default AccountSettingsBtn;
