"use client";

import React, { useEffect } from "react";
import "../../src/api/ws";
import { ProvideWS } from "../../src/contexts/WSContext";
import { useRouter } from "next/navigation";
import { useAuth } from "../../src/contexts/authContext";
import Loader from "../../src/components/common/Loader";

export default function Template({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isChecking } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated && !isChecking) {
      router.push("/login");
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, isChecking]);

  if (isChecking) {
    return <Loader />;
  }

  return <ProvideWS>{children}</ProvideWS>;
}
