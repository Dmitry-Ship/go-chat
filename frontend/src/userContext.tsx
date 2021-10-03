import React from "react";

type UserContextType = {
  id: string | null;
};

const user: UserContextType = {
  id: null,
};

export const UserContext = React.createContext(user);
