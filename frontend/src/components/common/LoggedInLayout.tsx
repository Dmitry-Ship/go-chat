import React, { useEffect } from "react";
import { ProvideWS } from "../../contexts/WSContext";
import { useRouter } from "next/router";
import { useAuth } from "../../contexts/authContext";
import Loader from "./Loader";

const LoggedInLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const auth = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!auth.isAuthenticated && !auth.isChecking) {
      router.push("/login");
    }
  }, [auth.isAuthenticated, auth.isChecking, router]);

  if (auth.isChecking) {
    return <Loader />;
  }

  return <ProvideWS isEnabled={true}>{children}</ProvideWS>;
};

export default LoggedInLayout;
