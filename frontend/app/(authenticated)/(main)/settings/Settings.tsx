"use client";
import React from "react";
import styles from "./Settings.module.css";
import Avatar from "../../../../src/components/common/Avatar";
import { useAuth } from "../../../../src/contexts/authContext";

function Settings() {
  const { user, logout } = useAuth();

  return (
    <>
      <header className={`header`}>
        <h2>Settings</h2>
      </header>
      <div className={styles.settingsPage}>
        <div className={styles.accountInfo}>
          <Avatar src={user?.avatar || ""} size={"lg"} />
          <h3>{user?.name}</h3>
        </div>
        <button onClick={logout} className={`btn`}>
          Logout
        </button>
      </div>
    </>
  );
}

export default Settings;
