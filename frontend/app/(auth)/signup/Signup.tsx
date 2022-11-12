"use client";

import React, { FormEvent, useState } from "react";
import styles from "../Login.module.css";
import { useAuth } from "../../../src/contexts/authContext";
import Link from "next/link";

function SignUp() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const { signup } = useAuth();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    signup(username, password);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        name="username"
        className="input"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />

      <input
        type="password"
        name="password"
        className="input m-top-1"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button
        disabled={!username || !password}
        type="submit"
        className="btn m-top-1"
      >
        SignUp
      </button>
      <Link href="/login" className={`m-top-1 ${styles.signUpLink}`}>
        I already have an account
      </Link>
    </form>
  );
}

export default SignUp;
