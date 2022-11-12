import React from "react";
import { ProvideAuth } from "../src/contexts/authContext";
import { ProvideAPI } from "../src/contexts/apiContext";
import ErrorAlert from "./ErrorAlert";

export default function Layouts({ children }: { children: React.ReactNode }) {
  return (
    <ProvideAuth>
      <ProvideAPI>
        <ErrorAlert />
        {children}
      </ProvideAPI>
    </ProvideAuth>
  );
}
