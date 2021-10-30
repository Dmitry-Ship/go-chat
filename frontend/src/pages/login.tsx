import React from "react";
import Login from "../components/authentication/Login";
import LoggedOutLayout from "../components/common/LoggedOutLayout";

function login() {
  return (
    <LoggedOutLayout>
      <Login />
    </LoggedOutLayout>
  );
}

export default login;
