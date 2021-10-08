import React from "react";
import ReactDOM from "react-dom";
import "./normalize.css";
import "./index.css";
import App from "./components/App";
import "./api/ws";

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root")
);
