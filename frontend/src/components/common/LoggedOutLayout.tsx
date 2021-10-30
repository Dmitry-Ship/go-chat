import React, { useEffect } from "react";
import { useRouter } from "next/router";
import { useAuth } from "../../contexts/authContext";
import Loader from "./Loader";

const LoggedOutLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const auth = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (auth.isAuthenticated && !auth.isChecking) {
      router.push("/");
    }
  }, [auth.isAuthenticated, auth.isChecking, router]);

  if (auth.isChecking) {
    return <Loader />;
  }

  return <>{children}</>;
};

export default LoggedOutLayout;
