"use client";
import React, { useEffect } from "react";
import styles from "./Login.module.css";
import { useRouter } from "next/navigation";
import { useAuth } from "../../src/contexts/authContext";
import { Loader } from "../../src/components/common/Loader";

function AuthLayout({ children }: { children: React.ReactNode }) {
  const { user, isChecking } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (user?.id) {
      router.push("/main");
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user?.id, isChecking]);

  return (
    <>
      <header className={`header`}>
        <h2>Go Chat</h2>
      </header>

      <section className={styles.wrap}>
        {isChecking || user ? <Loader /> : children}
      </section>
    </>
  );
}

export default AuthLayout;
