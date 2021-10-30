import React from "react";
import LoggedInLayout from "../components/common/LoggedInLayout";
import Navigation from "../components/Navigation";
import Rooms from "../components/rooms/Rooms";

function Index() {
  return (
    <LoggedInLayout>
      <Rooms />
      <Navigation />
    </LoggedInLayout>
  );
}

export default Index;
