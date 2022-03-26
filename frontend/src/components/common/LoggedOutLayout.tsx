import React, { useEffect } from "react";
import { useRouter } from "next/router";
import { useAuth } from "../../contexts/authContext";
import Loader from "./Loader";

const LoggedOutLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { isAuthenticated, isChecking } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isAuthenticated && !isChecking) {
      router.push("/");
    }
  }, [isAuthenticated, isChecking, router]);

  if (isChecking) {
    return <Loader />;
  }

  return <>{children}</>;
};

export default LoggedOutLayout;
