import React, { useEffect } from "react";
import styles from "./Layout.module.css";
import { ProvideAuth } from "../contexts/authContext";
import AuthLayout from "./common/AuthLayout";

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  useEffect(() => {
    const appHeight = () => {
      const doc = document.documentElement;
      doc.style.setProperty("--vh", `${window.innerHeight}px / 100`);
    };
    window.addEventListener("resize", appHeight);
    appHeight();
  }, []);

  return (
    <ProvideAuth>
      <AuthLayout>
        <div className={styles.app}>{children}</div>
      </AuthLayout>
    </ProvideAuth>
  );
};

export default Layout;
