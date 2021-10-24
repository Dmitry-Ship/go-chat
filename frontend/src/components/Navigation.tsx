import React from "react";
import { useAuth } from "../authContext";
import { Link, useRouteMatch } from "react-router-dom";
import AccountSettingsBtn from "./AccountSettingsBtn";

const Navigation = () => {
  const auth = useAuth();
  let match = useRouteMatch("/room/:roomId");

  if (!auth.isAuthenticated || match) {
    return null;
  }

  return (
    <div className="controls-for-scrollable">
      <Link to="/rooms" className="navBtn">
        💬
      </Link>
      <Link to="/people" className="navBtn">
        👥
      </Link>
      <AccountSettingsBtn />
    </div>
  );
};

export default Navigation;
