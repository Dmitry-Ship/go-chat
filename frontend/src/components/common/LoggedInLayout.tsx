import React, { useEffect } from "react";
import { ProvideWS } from "../../contexts/WSContext";
import { useRouter } from "next/router";
import { useAuth } from "../../contexts/authContext";
import Loader from "./Loader";

const LoggedInLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { isAuthenticated, isChecking } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated && !isChecking) {
      router.push("/login");
    }
  }, [isAuthenticated, isChecking, router]);

  if (isChecking) {
    return <Loader />;
  }

  return <ProvideWS isEnabled={true}>{children}</ProvideWS>;
};

export default LoggedInLayout;
