import React from "react";
import LoggedInLayout from "../components/common/LoggedInLayout";
import Navigation from "../components/Navigation";
import SettingsPage from "../components/settings/Settings";

function Settings() {
  return (
    <LoggedInLayout>
      <SettingsPage />
      <Navigation />
    </LoggedInLayout>
  );
}

export default Settings;
