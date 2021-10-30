import "../styles/index.css";
import "../styles/normalize.css";
import "../api/ws";
import React from "react";
import Layout from "../components/Layout";
import { AppProps } from "next/app";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <Layout>
      <Component {...pageProps} />
    </Layout>
  );
}
