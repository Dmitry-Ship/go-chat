import React from "react";
import Link from "next/link";
import AccountSettingsBtn from "./AccountSettingsBtn";

const Navigation = () => {
  return (
    <div className="controls-for-scrollable">
      <Link href="/">
        <a className="navBtn">💬</a>
      </Link>
      <Link href="/people">
        <a className="navBtn">👥</a>
      </Link>
      <AccountSettingsBtn />
    </div>
  );
};

export default Navigation;
