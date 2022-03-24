import React from "react";
import styles from "./Settings.module.css";
import Avatar from "../common/Avatar";
import { useAuth } from "../../contexts/authContext";

function Settings() {
  const { user, logout } = useAuth();

  return (
    <div className={styles.settingsPage}>
      <div className={styles.accountInfo}>
        <Avatar src={user?.avatar || ""} size={100} />
        <h3>{user?.name}</h3>
      </div>
      <button onClick={logout} className={`btn`}>
        Logout
      </button>
    </div>
  );
}

export default Settings;
