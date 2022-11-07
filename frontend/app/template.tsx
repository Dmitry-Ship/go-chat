"use client";

import React, { useEffect } from "react";
import { ProvideAuth } from "../src/contexts/authContext";
import { ProvideAPI } from "../src/contexts/apiContext";
import ErrorAlert from "./ErrorAlert";

export default function Template({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    const appHeight = () => {
      const doc = document.documentElement;
      doc.style.setProperty("--vh", `${window.innerHeight}px / 100`);
    };
    window.addEventListener("resize", appHeight);
    appHeight();
  }, []);

  return (
    <ProvideAPI>
      <ErrorAlert />
      <ProvideAuth>{children}</ProvideAuth>
    </ProvideAPI>
  );
}
