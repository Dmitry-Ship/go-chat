import React, { useEffect } from "react";
import { ProvideWS } from "../../contexts/WSContext";
import { useRouter } from "next/router";
import { useAuth } from "../../contexts/authContext";
import Loader from "./Loader";

const AuthLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, isChecking } = useAuth();
  const router = useRouter();

  useEffect(() => {
    router.push(isAuthenticated ? "/" : "/login");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  if (isChecking) {
    return <Loader />;
  }

  return isAuthenticated ? <ProvideWS>{children}</ProvideWS> : <>{children}</>;
};

export default AuthLayout;
