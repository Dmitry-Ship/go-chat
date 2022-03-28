import React from "react";
import LoggedInLayout from "../components/common/LoggedInLayout";
import ContactsList from "../components/contacts/ContactsList";
import Navigation from "../components/Navigation";

function Settings() {
  return (
    <LoggedInLayout>
      <ContactsList />
      <Navigation />
    </LoggedInLayout>
  );
}

export default Settings;
