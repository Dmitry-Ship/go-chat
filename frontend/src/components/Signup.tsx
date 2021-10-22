import React, { FormEvent } from "react";
import styles from "./Login.module.css";
import { useAuth } from "../authContext";

function SignUp() {
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const { signup } = useAuth();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    signup(username, password);
  };

  return (
    <>
      <header className={`header`}>
        <h2>Sign Up</h2>
      </header>
      <section className={styles.wrap}>
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
        </form>
      </section>
    </>
  );
}

export default SignUp;
