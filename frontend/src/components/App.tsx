import React from "react";
import styles from "./App.module.css";
import Chat from "./ChatRoom/Chat";
import { BrowserRouter as Router, Switch } from "react-router-dom";
import Rooms from "./Rooms";
import { useWS } from "../api/hooks";
import Loader from "./Loader";
import Login from "./Login";
import PrivateRoute from "./PrivateRoute";
import { ProvideAuth, useAuth } from "../authContext";
import AuthRoute from "./AuthRoute";
import SignUp from "./Signup";

function App() {
  return (
    <ProvideAuth>
      <div className={styles.app}>
        <Routes />
      </div>
    </ProvideAuth>
  );
}

const Routes = () => {
  const auth = useAuth();
  const { status } = useWS();

  if (auth.isChecking || status === "connecting") {
    return <Loader />;
  }

  return (
    <Router>
      <Switch>
        <PrivateRoute path="/room/:roomId">
          <Chat />
        </PrivateRoute>
        <AuthRoute path="/login">
          <Login />
        </AuthRoute>
        <AuthRoute path="/signup">
          <SignUp />
        </AuthRoute>
        <PrivateRoute path="/">
          <Rooms />
        </PrivateRoute>
      </Switch>
    </Router>
  );
};

export default App;
