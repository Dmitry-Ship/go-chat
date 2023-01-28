"use client";

import React, { useEffect } from "react";
import "../../src/api/ws";
import { ProvideWS } from "../../src/contexts/WSContext";
import { useRouter } from "next/navigation";
import { useAuth } from "../../src/contexts/authContext";
import Loader from "../../src/components/common/Loader";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, isChecking } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!(user || isChecking)) {
      router.push("/login");
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user, isChecking]);

  useEffect(() => {
    const appHeight = () => {
      const doc = document.documentElement;
      doc.style.setProperty("--vh", `${window.innerHeight}px / 100`);
    };
    window.addEventListener("resize", appHeight);
    appHeight();
  }, []);

  if (isChecking || !user) {
    return <Loader />;
  }

  return <ProvideWS>{children}</ProvideWS>;
}
