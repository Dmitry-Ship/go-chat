import React from "react";
import SignUp from "../components/authentication/Signup";
import LoggedOutLayout from "../components/common/LoggedOutLayout";

function signup() {
  return (
    <LoggedOutLayout>
      <SignUp />
    </LoggedOutLayout>
  );
}

export default signup;
