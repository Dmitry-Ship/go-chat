import React, { useState } from "react";
import styles from "./AccountSettingsBtn.module.css";
import SlideIn from "./SlideIn";
import { useAuth } from "../authContext";

const AccountSettingsBtn: React.FC = () => {
  const [isEditing, setIsEditing] = useState(false);
  const auth = useAuth();

  return (
    <>
      <button
        onClick={() => setIsEditing(true)}
        className={styles.accountSettingsBtn}
      >
        ⚙️
      </button>
      <SlideIn onClose={() => setIsEditing(false)} isOpen={isEditing}>
        <>
          <button onClick={auth.logout} className={`btn`}>
            Logout
          </button>
        </>
      </SlideIn>
    </>
  );
};

export default AccountSettingsBtn;
