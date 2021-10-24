import React from "react";
import { Redirect, Route } from "react-router-dom";
import { useAuth } from "../../authContext";

const AuthRoute: React.FC<{
  children: React.ReactNode;
  [key: string]: any;
}> = ({ children, ...rest }) => {
  const auth = useAuth();

  return (
    <Route
      {...rest}
      render={({ location }) =>
        !auth.isAuthenticated ? (
          children
        ) : (
          <Redirect
            to={{
              pathname: "/",
              state: { from: location },
            }}
          />
        )
      }
    />
  );
};

export default AuthRoute;
