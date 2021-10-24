import React, { useEffect } from "react";
import { useAuth } from "../authContext";
import { ProvideWS, useWS } from "../WSContext";
import { BrowserRouter as Router, Switch } from "react-router-dom";
import Loader from "./common/Loader";
import AuthRoute from "./common/AuthRoute";
import Login from "./authentication/Login";
import SignUp from "./authentication/Signup";
import Chat from "./ChatRoom/Chat";
import PrivateRoute from "./common/PrivateRoute";
import Rooms from "./Rooms";
import Navigation from "./Navigation";

const Routes = () => {
  const auth = useAuth();

  return (
    <ProvideWS isEnabled={auth.isAuthenticated}>
      <Router>
        <Switch>
          <AuthRoute path="/login">
            <Login />
          </AuthRoute>
          <AuthRoute path="/signup">
            <SignUp />
          </AuthRoute>
          <PrivateRoute path="/room/:roomId">
            <Chat />
          </PrivateRoute>
          <PrivateRoute path="/">
            <Rooms />
          </PrivateRoute>
        </Switch>

        <Navigation />
      </Router>
    </ProvideWS>
  );
};

export default Routes;
